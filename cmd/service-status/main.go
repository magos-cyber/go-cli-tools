package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	checkServiceName := checkCmd.String("service", "", "Service name")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCmd.Bool("all", false, "List all services")

	if len(os.Args) < 2 {
		fmt.Println("Usage: service-status <command> [arguments]")
		fmt.Println("Commands: check, list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkCmd.Parse(os.Args[2:])
		if *checkServiceName == "" {
			fmt.Println("Error: --service is required")
			os.Exit(1)
		}
		doCheckService(*checkServiceName)
	case "list":
		listCmd.Parse(os.Args[2:])
		listServices(*listAll)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func doCheckService(name string) {
	cmd := exec.Command("systemctl", "is-active", name)
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))

	if err != nil {
		fmt.Printf("Service %s: %s\n", name, status)
		os.Exit(1)
	}

	fmt.Printf("Service %s: %s\n", name, status)
}

func listServices(all bool) {
	args := []string{"list-units", "--type=service"}
	if !all {
		args = append(args, "--state=running")
	}
	cmd := exec.Command("systemctl", args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Println(line)
		}
	}
}
