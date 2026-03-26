//go:build linux
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s run <command> [args...]\n", os.Args[0])
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

	// Re-execute itself with 'child' as the first argument
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Configure namespaces: PID, UTS (hostname), Mount (filesystem)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// child is the second stage running inside the new namespaces
func child() {
	fmt.Printf("Running Stage 2 (PID: %d in container)\n", os.Getpid())

	// Set a new hostname for the UTS namespace
	must(syscall.Sethostname([]byte("container-runtime")))

	// Mount /proc to isolate PID namespace visibility
	// MS_NOEXEC: Do not allow programs to be executed from this filesystem
	// MS_NOSUID: Do not honor set-user-ID or set-group-ID bits
	// MS_NODEV: Do not allow access to devices (special files)
	defaultMountFlags := uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV)
	must(syscall.Mount("proc", "/proc", "proc", defaultMountFlags, ""))

	// Execute the final user command, replacing this process
	cmdPath, err := exec.LookPath(os.Args[2])
	must(err)

	must(syscall.Exec(cmdPath, os.Args[2:], os.Environ()))
}

func must(err error) {
	if err != nil {
		fmt.Printf("Fatal Error: %v\n", err)
		os.Exit(1)
	}
}
