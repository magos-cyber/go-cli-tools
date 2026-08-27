package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployManifest := deployCmd.String("manifest", "", "Path to Kubernetes manifest file")
	deployNamespace := deployCmd.String("namespace", "default", "Target namespace")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listNamespace := listCmd.String("namespace", "default", "Namespace to list resources")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteResource := deleteCmd.String("resource", "", "Resource to delete (e.g., deployment/nginx)")
	deleteNamespace := deleteCmd.String("namespace", "default", "Namespace of the resource")

	if len(os.Args) < 2 {
		fmt.Println("Usage: k8s-deploy <command> [arguments]")
		fmt.Println("Commands: deploy, list, delete")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		deployCmd.Parse(os.Args[2:])
		if *deployManifest == "" {
			fmt.Println("Error: --manifest is required")
			os.Exit(1)
		}
		fmt.Printf("Deploying manifest '%s' to namespace '%s'...\n", *deployManifest, *deployNamespace)
		
		args := []string{"apply", "-f", *deployManifest, "-n", *deployNamespace}
		cmd := exec.Command("kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to deploy manifest: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Manifest deployed successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Printf("Listing resources in namespace '%s'...\n", *listNamespace)
		
		args := []string{"get", "all", "-n", *listNamespace}
		cmd := exec.Command("kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to list resources: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteResource == "" {
			fmt.Println("Error: --resource is required")
			os.Exit(1)
		}
		fmt.Printf("Deleting resource '%s' from namespace '%s'...\n", *deleteResource, *deleteNamespace)
		
		args := strings.Split(*deleteResource, "/")
		if len(args) != 2 {
			fmt.Println("Error: --resource must be in format type/name")
			os.Exit(1)
		}
		
		cmdArgs := []string{"delete", args[0], args[1], "-n", *deleteNamespace}
		cmd := exec.Command("kubectl", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to delete resource: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Resource deleted successfully.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}