package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	path := flag.String"path", "/", "Path to check")
	flag.Parse()

	out, err := exec.Command("df", "-h", *path).Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "%") {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				fmt.Printf("Filesystem: %s\n", fields[0])
				fmt.Printf("Size: %s\n", fields[1])
				fmt.Printf("Used: %s\n", fields[2])
				fmt.Printf("Available: %s\n", fields[3])
				fmt.Printf("Use%%: %s\n", fields[4])
				fmt.Printf("Mounted: %s\n", fields[5])
			}
		}
	}
}
