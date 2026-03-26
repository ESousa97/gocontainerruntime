//go:build linux

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	// cacheDir is the local directory where the rootfs will be extracted.
	cacheDir = "./cache/alpine_rootfs"
	// alpineURL is the official Alpine Linux minirootfs download URL.
	alpineURL = "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.1-x86_64.tar.gz"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gocontainer",
		Short: "A minimal container runtime implemented in Go",
	}

	// pullCmd handles downloading and extracting the rootfs.
	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Download a basic Alpine rootfs for the container",
		Run: func(cmd *cobra.Command, args []string) {
			pull()
		},
	}

	// runCmd is the main entry point to start a container.
	var runCmd = &cobra.Command{
		Use:   "run [rootfs_path] [command] [args...]",
		Short: "Run a command inside a new container",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var rootfs, userCommand string
			var userArgs []string

			// If only 1 arg, use cached image and the arg is the command
			if len(args) == 1 {
				if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
					fmt.Println("No rootfs provided and cache is empty. Running 'pull' first...")
					pull()
				}
				rootfs = cacheDir
				userCommand = args[0]
				userArgs = args
			} else {
				rootfs = args[0]
				userCommand = args[1]
				userArgs = args[1:]
			}

			run(rootfs, userCommand, userArgs)
		},
	}

	// childCmd is an internal command used for re-execution in a new namespace.
	var childCmd = &cobra.Command{
		Use:    "child",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] is rootfs, args[1] is command
			child(args[0], args[1], args[1:])
		},
	}

	rootCmd.AddCommand(pullCmd, runCmd, childCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Execution Error: %v\n", err)
		os.Exit(1)
	}
}

// pull downloads and extracts the Alpine Linux minirootfs to the cache directory.
func pull() {
	fmt.Printf("Downloading Alpine rootfs from %s...\n", alpineURL)

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		must(err)
	}

	resp, err := http.Get(alpineURL)
	must(err)
	defer resp.Body.Close()

	uncompressed, err := gzip.NewReader(resp.Body)
	must(err)

	tr := tar.NewReader(uncompressed)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		must(err)

		target := filepath.Join(cacheDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			must(err)
			io.Copy(f, tr)
			f.Close()
		}
	}
	fmt.Printf("Download and extraction complete. Cache ready at: %s\n", cacheDir)
}

// run performs the first stage of container creation: setting up Namespaces and Cgroups.
// It then re-executes itself with the 'child' command inside the new isolation layers.
func run(rootfs, userCommand string, userArgs []string) {
	fmt.Printf("Running Stage 1 (PID: %d)\n", os.Getpid())

	// Re-execute itself with 'child' as the first argument to run inside the namespace.
	args := append([]string{"child", rootfs, userCommand}, userArgs...)
	cmd := exec.Command("/proc/self/exe", args...)

	// Configure namespaces: PID (processes), UTS (hostname), Mount (fs), Network (ip).
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Cgroups setup (using v1 paths)
	cg := "/sys/fs/cgroup"
	memCg := cg + "/memory/gocontainer"
	cpuCg := cg + "/cpu/gocontainer"

	must(os.MkdirAll(memCg, 0755))
	must(os.MkdirAll(cpuCg, 0755))

	defer func() {
		os.Remove(memCg)
		os.Remove(cpuCg)
	}()

	// Apply resource limits: 100MB memory and 512 CPU shares.
	must(os.WriteFile(memCg+"/memory.limit_in_bytes", []byte("104857600"), 0644))
	must(os.WriteFile(cpuCg+"/cpu.shares", []byte("512"), 0644))

	must(cmd.Start())

	// Network Setup: creates a veth pair and moves one end into the container's netns.
	pid := cmd.Process.Pid
	exec.Command("ip", "link", "add", "veth-host", "type", "veth", "peer", "name", "veth-child").Run()
	exec.Command("ip", "link", "set", "veth-child", "netns", fmt.Sprintf("%d", pid)).Run()
	exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", "veth-host").Run()
	exec.Command("ip", "link", "set", "veth-host", "up").Run()

	// Configure networking inside the container using nsenter.
	nsenter := []string{"nsenter", "-t", fmt.Sprintf("%d", pid), "-n"}
	exec.Command(nsenter[0], append(nsenter[1:], "ip", "addr", "add", "10.0.0.2/24", "dev", "veth-child")...).Run()
	exec.Command(nsenter[0], append(nsenter[1:], "ip", "link", "set", "veth-child", "up")...).Run()
	exec.Command(nsenter[0], append(nsenter[1:], "ip", "link", "set", "lo", "up")...).Run()

	defer exec.Command("ip", "link", "delete", "veth-host").Run()

	// Write PID to Cgroups to start enforcing limits.
	pidStr := []byte(fmt.Sprintf("%d", pid))
	must(os.WriteFile(memCg+"/cgroup.procs", pidStr, 0644))
	must(os.WriteFile(cpuCg+"/cgroup.procs", pidStr, 0644))

	must(cmd.Wait())
}

// child performs the second stage of container setup inside the isolated namespaces.
// It sets the hostname, chroots into the filesystem, and mounts /proc.
func child(rootfs, userCommand string, userArgs []string) {
	fmt.Printf("Running Stage 2 (PID: %d in container)\n", os.Getpid())

	must(syscall.Sethostname([]byte("gocontainer")))
	must(syscall.Chroot(rootfs))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	// Find the command path inside the isolated rootfs.
	cmdPath, err := exec.LookPath(userCommand)
	if err != nil {
		fmt.Printf("Execution path error: %v\n", err)
		os.Exit(1)
	}

	// Replace the current process with the user's command.
	must(syscall.Exec(cmdPath, userArgs, os.Environ()))
}

// must is a helper function that checks for errors and terminates the process if one occurs.
func must(err error) {
	if err != nil {
		fmt.Printf("Fatal Runtime Error: %v\n", err)
		os.Exit(1)
	}
}
