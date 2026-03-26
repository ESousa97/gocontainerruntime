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

	// Configure namespaces: PID, UTS (hostname), Mount (filesystem)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
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

	// Ensure cleanup
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
