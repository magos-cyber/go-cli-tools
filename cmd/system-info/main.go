package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("=== System Info ===")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	out, _ := exec.Command("free", "-h").Output()
	fmt.Printf("Memory:\n%s\n", string(out))
	out, _ = exec.Command("df", "-h").Output()
	fmt.Printf("Disk:\n%s\n", string(out))
	_ = strings.Fields("")
}
