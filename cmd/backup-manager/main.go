package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createSource := createCmd.String("source", "", "Source directory")
	createDest := createCmd.String("dest", "", "Destination directory")
	createName := createCmd.String("name", "", "Backup name (optional)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listDir := listCmd.String("dir", ".", "Directory to list backups")

	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	restoreFile := restoreCmd.String("file", "", "Backup file to restore")
	restoreDest := restoreCmd.String("dest", "", "Restore destination")

	if len(os.Args) < 2 {
		fmt.Println("Usage: backup-manager <command> [arguments]")
		fmt.Println("Commands: create, list, restore")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createSource == "" || *createDest == "" {
			fmt.Println("Error: --source and --dest are required")
			os.Exit(1)
		}
		createBackup(*createSource, *createDest, *createName)
	case "list":
		listCmd.Parse(os.Args[2:])
		listBackups(*listDir)
	case "restore":
		restoreCmd.Parse(os.Args[2:])
		if *restoreFile == "" || *restoreDest == "" {
			fmt.Println("Error: --file and --dest are required")
			os.Exit(1)
		}
		restoreBackup(*restoreFile, *restoreDest)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func createBackup(source, dest, name string) {
	if name == "" {
		name = fmt.Sprintf("backup-%s", time.Now().Format("20060102-150405"))
	}
	
	archive := filepath.Join(dest, name+".tar.gz")
	
	fmt.Printf("Creating backup: %s\n", archive)
	fmt.Printf("Source: %s\n", source)
	
	cmd := exec.Command("tar", "-czf", archive, "-C", filepath.Dir(source), filepath.Base(source))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Backup created successfully")
}

func listBackups(dir string) {
	fmt.Println("=== Available Backups ===")
	
	pattern := filepath.Join(dir, "backup-*.tar.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("Error listing backups: %v\n", err)
		os.Exit(1)
	}
	
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		fmt.Printf("  %s (%d MB)\n", filepath.Base(match), info.Size()/1024/1024)
	}
}

func restoreBackup(file, dest string) {
	fmt.Printf("Restoring: %s\n", file)
	fmt.Printf("Destination: %s\n", dest)
	
	cmd := exec.Command("tar", "-xzf", file, "-C", dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Restore failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Restore completed successfully")
}
