//go:build linux
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s run <rootfs_path> <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Invalid command")
	}
}

// run is the first stage that prepares the namespaces
func run() {
	fmt.Printf("Running Stage 1 (PID: %d)\n", os.Getpid())

	// Re-execute itself with 'child' as the first argument, passing rootfs and command
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Configure namespaces: PID, UTS (hostname), Mount (filesystem), Network
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up Cgroups
	cg := "/sys/fs/cgroup"
	memCg := cg + "/memory/gocontainer"
	cpuCg := cg + "/cpu/gocontainer"

	// Create cgroup directories
	must(os.MkdirAll(memCg, 0755))
	must(os.MkdirAll(cpuCg, 0755))

	// Ensure cgroup cleanup
	defer func() {
		os.Remove(memCg)
		os.Remove(cpuCg)
	}()

	// Set memory limit: 100MB
	must(os.WriteFile(memCg+"/memory.limit_in_bytes", []byte("104857600"), 0644))
	// Set CPU shares
	must(os.WriteFile(cpuCg+"/cpu.shares", []byte("512"), 0644))

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting child process: %v\n", err)
		os.Exit(1)
	}

	// Network Setup
	pid := cmd.Process.Pid
	fmt.Printf("Setting up network for PID %d\n", pid)

	// 1. Create veth pair
	must(exec.Command("ip", "link", "add", "veth-host", "type", "veth", "peer", "name", "veth-child").Run())
	
	// 2. Move veth-child to the container's network namespace
	must(exec.Command("ip", "link", "set", "veth-child", "netns", fmt.Sprintf("%d", pid)).Run())

	// 3. Configure the host side
	must(exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", "veth-host").Run())
	must(exec.Command("ip", "link", "set", "veth-host", "up").Run())

	// 4. Configure the child side (inside the namespace)
	// We use 'nsenter' to run commands inside the child's network namespace
	nsenter := []string{"nsenter", "-t", fmt.Sprintf("%d", pid), "-n"}
	must(exec.Command(nsenter[0], append(nsenter[1:], "ip", "addr", "add", "10.0.0.2/24", "dev", "veth-child")...).Run())
	must(exec.Command(nsenter[0], append(nsenter[1:], "ip", "link", "set", "veth-child", "up")...).Run())
	must(exec.Command(nsenter[0], append(nsenter[1:], "ip", "link", "set", "lo", "up")...).Run())

	// Ensure network cleanup on exit
	defer func() {
		exec.Command("ip", "link", "delete", "veth-host").Run()
	}()

	// Write the child process PID to cgroup.procs
	pidStr := []byte(fmt.Sprintf("%d", cmd.Process.Pid))
	must(os.WriteFile(memCg+"/cgroup.procs", pidStr, 0644))
	must(os.WriteFile(cpuCg+"/cgroup.procs", pidStr, 0644))

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for child process: %v\n", err)
		os.Exit(1)
	}
}

// child is the second stage running inside the new namespaces
func child() {
	fmt.Printf("Running Stage 2 (PID: %d in container)\n", os.Getpid())

	rootfs := os.Args[2]
	userCommand := os.Args[3]
	userArgs := os.Args[3:]

	// Set a new hostname for the UTS namespace
	must(syscall.Sethostname([]byte("container-runtime")))

	// 1. Isolate the filesystem: Chroot to the provided rootfs path
	must(syscall.Chroot(rootfs))

	// 2. Change directory to the new root
	must(os.Chdir("/"))

	// 3. Mount /proc inside the new root to isolate PID visibility
	// This must happen after chroot so it is mounted in the container's /proc
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	// Execute the final user command, replacing this process
	// Since we are inside the chroot, we look for the command relative to the new root
	cmdPath, err := exec.LookPath(userCommand)
	if err != nil {
		fmt.Printf("Error finding command '%s' inside rootfs: %v\n", userCommand, err)
		os.Exit(1)
	}

	must(syscall.Exec(cmdPath, userArgs, os.Environ()))
}

func must(err error) {
	if err != nil {
		fmt.Printf("Fatal Error: %v\n", err)
		os.Exit(1)
	}
}
