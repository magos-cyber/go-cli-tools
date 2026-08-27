package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployStack := deployCmd.String("stack", "", "Name of the stack to deploy")
	deployEnvFile := deployCmd.String("env-file", "", "Path to .env file")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateAll := updateCmd.Bool("all", false, "Update all stacks")
	updateStack := updateCmd.String("stack", "", "Update specific stack")

	logsCmd := flag.NewFlagSet("logs", flag.ExitOnError)
	logsStack := logsCmd.String("stack", "", "Stack name for logs")
	logsTail := logsCmd.Int("tail", 100, "Number of tail lines")

	if len(os.Args) < 2 {
		fmt.Println("Usage: docker-compose <command> [arguments]")
		fmt.Println("Commands: deploy, list, update, logs")
		os.Exit(1)
	}

	stacksDir := os.Getenv("DOCKER_STACKS_DIR")
	if stacksDir == "" {
		stacksDir = "/opt/stacks"
	}

	switch os.Args[1] {
	case "deploy":
		deployCmd.Parse(os.Args[2:])
		if *deployStack == "" {
			fmt.Println("Error: --stack is required")
			os.Exit(1)
		}
		stackPath := filepath.Join(stacksDir, *deployStack)
		fmt.Printf("Deploying stack '%s' from %s...\n", *deployStack, stackPath)

		args := []string{"compose", "up", "-d"}
		if *deployEnvFile != "" {
			args = append([]string{"--env-file", *deployEnvFile}, args...)
		}

		cmd := exec.Command("docker", args...)
		cmd.Dir = stackPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to deploy stack: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Stack deployed successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Println("=== Active Docker Compose Stacks ===")
		
		files, err := os.ReadDir(stacksDir)
		if err != nil {
			fmt.Printf("Could not read stacks directory (%s): %v\n", stacksDir, err)
			os.Exit(1)
		}

		for _, f := range files {
			if f.IsDir() {
				stackName := f.Name()
				stackPath := filepath.Join(stacksDir, stackName)
				// check if docker-compose.yml or compose.yml exists
				if fileExists(filepath.Join(stackPath, "docker-compose.yml")) || fileExists(filepath.Join(stackPath, "compose.yml")) {
					fmt.Printf("- %s (%s)\n", stackName, stackPath)
				}
			}
		}

	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateAll {
			fmt.Println("Updating all stacks...")
			files, err := os.ReadDir(stacksDir)
			if err != nil {
				fmt.Printf("Could not read stacks directory: %v\n", err)
				os.Exit(1)
			}
			for _, f := range files {
				if f.IsDir() {
					updateStackDir(filepath.Join(stacksDir, f.Name()))
				}
			}
		} else if *updateStack != "" {
			updateStackDir(filepath.Join(stacksDir, *updateStack))
		} else {
			fmt.Println("Error: specify --stack <name> or --all")
			os.Exit(1)
		}

	case "logs":
		logsCmd.Parse(os.Args[2:])
		if *logsStack == "" {
			fmt.Println("Error: --stack is required")
			os.Exit(1)
		}
		stackPath := filepath.Join(stacksDir, *logsStack)
		args := []string{"compose", "logs", "--tail", fmt.Sprintf("%d", *logsTail), "-f"}

		cmd := exec.Command("docker", args...)
		cmd.Dir = stackPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error fetching logs: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func updateStackDir(dir string) {
	if !fileExists(filepath.Join(dir, "docker-compose.yml")) && !fileExists(filepath.Join(dir, "compose.yml")) {
		return
	}
	fmt.Printf("Updating stack in %s...\n", dir)
	cmdPull := exec.Command("docker", "compose", "pull")
	cmdPull.Dir = dir
	cmdPull.Stdout = os.Stdout
	cmdPull.Stderr = os.Stderr
	_ = cmdPull.Run()

	cmdUp := exec.Command("docker", "compose", "up", "-d")
	cmdUp.Dir = dir
	cmdUp.Stdout = os.Stdout
	cmdUp.Stderr = os.Stderr
	_ = cmdUp.Run()
	strings.TrimSpace("") //noop
}
