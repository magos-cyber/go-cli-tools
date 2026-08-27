package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	installName := installCmd.String("name", "", "Release name")
	installChart := installCmd.String("chart", "", "Chart name or path")
	installNamespace := installCmd.String("namespace", "default", "Target namespace")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	upgradeCmd := flag.NewFlagSet("upgrade", flag.ExitOnError)
	upgradeName := upgradeCmd.String("name", "", "Release name")
	upgradeChart := upgradeCmd.String("chart", "", "Chart name or path")
	upgradeNamespace := upgradeCmd.String("namespace", "default", "Target namespace")

	uninstallCmd := flag.NewFlagSet("uninstall", flag.ExitOnError)
	uninstallName := uninstallCmd.String("name", "", "Release name")
	uninstallNamespace := uninstallCmd.String("namespace", "default", "Target namespace")

	if len(os.Args) < 2 {
		fmt.Println("Usage: helm-cli <command> [arguments]")
		fmt.Println("Commands: install, list, upgrade, uninstall")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		if *installName == "" || *installChart == "" {
			fmt.Println("Error: --name and --chart are required")
			os.Exit(1)
		}
		fmt.Printf("Installing chart '%s' as release '%s' in namespace '%s'...\n", *installChart, *installName, *installNamespace)
		
		args := []string{"install", *installName, *installChart, "--namespace", *installNamespace}
		cmd := exec.Command("helm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to install chart: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Chart installed successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Println("Listing Helm releases...")
		
		args := []string{"list", "--all-namespaces"}
		cmd := exec.Command("helm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to list releases: %v\n", err)
			os.Exit(1)
		}

	case "upgrade":
		upgradeCmd.Parse(os.Args[2:])
		if *upgradeName == "" || *upgradeChart == "" {
			fmt.Println("Error: --name and --chart are required")
			os.Exit(1)
		}
		fmt.Printf("Upgrading release '%s' with chart '%s' in namespace '%s'...\n", *upgradeName, *upgradeChart, *upgradeNamespace)
		
		args := []string{"upgrade", *upgradeName, *upgradeChart, "--namespace", *upgradeNamespace}
		cmd := exec.Command("helm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to upgrade release: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Release upgraded successfully.")

	case "uninstall":
		uninstallCmd.Parse(os.Args[2:])
		if *uninstallName == "" {
			fmt.Println("Error: --name is required")
			os.Exit(1)
		}
		fmt.Printf("Uninstalling release '%s' from namespace '%s'...\n", *uninstallName, *uninstallNamespace)
		
		args := []string{"uninstall", *uninstallName, "--namespace", *uninstallNamespace}
		cmd := exec.Command("helm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to uninstall release: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Release uninstalled successfully.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}