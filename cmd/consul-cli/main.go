package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	registerCmd := flag.NewFlagSet("register", flag.ExitOnError)
	registerName := registerCmd.String("name", "", "Service name")
	registerID := registerCmd.String("id", "", "Service ID")
	registerAddress := registerCmd.String("address", "", "Service address")
	registerPort := registerCmd.Int("port", 0, "Service port")

	deregisterCmd := flag.NewFlagSet("deregister", flag.ExitOnError)
	deregisterID := deregisterCmd.String("id", "", "Service ID to deregister")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: consul-cli <command> [arguments]")
		fmt.Println("Commands: register, deregister, list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "register":
		registerCmd.Parse(os.Args[2:])
		if *registerName == "" || *registerID == "" || *registerAddress == "" || *registerPort == 0 {
			fmt.Println("Error: --name, --id, --address, and --port are required")
			os.Exit(1)
		}
		fmt.Printf("Registering service '%s' with ID '%s'...\n", *registerName, *registerID)
		
		args := []string{"services", "register", "-name", *registerName, "-id", *registerID, "-address", *registerAddress, "-port", fmt.Sprintf("%d", *registerPort)}
		cmd := exec.Command("consul", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to register service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service registered successfully.")

	case "deregister":
		deregisterCmd.Parse(os.Args[2:])
		if *deregisterID == "" {
			fmt.Println("Error: --id is required")
			os.Exit(1)
		}
		fmt.Printf("Deregistering service with ID '%s'...\n", *deregisterID)
		
		args := []string{"services", "deregister", "-id", *deregisterID}
		cmd := exec.Command("consul", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to deregister service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service deregistered successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Println("Listing registered services...")
		
		args := []string{"services", "list"}
		cmd := exec.Command("consul", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to list services: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}