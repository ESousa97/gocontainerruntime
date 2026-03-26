//go:build windows
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("O Container Runtime em Go utiliza Namespaces do Linux (UTS, PID, Mount).")
	fmt.Println("Não é possível executá-lo nativamente no Windows.")
	fmt.Println("\nPara testar:")
	fmt.Println("1. Use o WSL2 (Windows Subsystem for Linux).")
	fmt.Println("2. No terminal do WSL, execute: go build -o runtime main.go")
	fmt.Println("3. Rode com privilégios de root: sudo ./runtime run /bin/bash")
	os.Exit(0)
}
