package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readPath := readCmd.String("path", "", "Path to secret")

	writeCmd := flag.NewFlagSet("write", flag.ExitOnError)
	writePath := writeCmd.String("path", "", "Path to write secret")
	writeData := writeCmd.String("data", "", "Data to write (key=value pairs)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listPath := listCmd.String("path", "", "Path to list secrets")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deletePath := deleteCmd.String("path", "", "Path to delete secret")

	if len(os.Args) < 2 {
		fmt.Println("Usage: vault-cli <command> [arguments]")
		fmt.Println("Commands: read, write, list, delete")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "read":
		readCmd.Parse(os.Args[2:])
		if *readPath == "" {
			fmt.Println("Error: --path is required")
			os.Exit(1)
		}
		fmt.Printf("Reading secret at path '%s'...\n", *readPath)
		
		args := []string{"kv", "get", "-format=json", *readPath}
		cmd := exec.Command("vault", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to read secret: %v\n", err)
			os.Exit(1)
		}

	case "write":
		writeCmd.Parse(os.Args[2:])
		if *writePath == "" || *writeData == "" {
			fmt.Println("Error: --path and --data are required")
			os.Exit(1)
		}
		fmt.Printf("Writing secret to path '%s'...\n", *writePath)
		
		args := []string{"kv", "put", *writePath}
		dataPairs := strings.Split(*writeData, ",")
		for _, pair := range dataPairs {
			args = append(args, strings.TrimSpace(pair))
		}
		
		cmd := exec.Command("vault", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to write secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Secret written successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		if *listPath == "" {
			fmt.Println("Error: --path is required")
			os.Exit(1)
		}
		fmt.Printf("Listing secrets at path '%s'...\n", *listPath)
		
		args := []string{"kv", "list", *listPath}
		cmd := exec.Command("vault", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to list secrets: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deletePath == "" {
			fmt.Println("Error: --path is required")
			os.Exit(1)
		}
		fmt.Printf("Deleting secret at path '%s'...\n", *deletePath)
		
		args := []string{"kv", "delete", *deletePath}
		cmd := exec.Command("vault", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to delete secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Secret deleted successfully.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}