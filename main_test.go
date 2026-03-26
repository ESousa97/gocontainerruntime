//go:build linux

package main

import (
	"os"
	"os/exec"
	"testing"
)

// TestRunCommand checks if the basic structure is sound
func TestRunCommand(t *testing.T) {
	// Only run this test in a Linux environment with root access
	if os.Getuid() != 0 {
		t.Skip("Skip test because it requires root privileges to create namespaces.")
	}

	// Try to execute the runtime with a simple command inside a container
	// This ensures the Stages 1 and 2 are talking correctly
	cmd := exec.Command("./runtime", "run", "echo", "hello")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed to run container: %v, output: %s", err, string(output))
	}

	if string(output) == "" {
		t.Errorf("No output from container")
	}
}
