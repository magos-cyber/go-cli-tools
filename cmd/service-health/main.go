package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HealthCheck struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Latency string `json:"latency"`
}

func main() {
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	checkURL := checkCmd.String("url", "", "URL to check")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listFile := listCmd.String("file", "", "File with URLs to check")

	if len(os.Args) < 2 {
		fmt.Println("Usage: service-health <command> [arguments]")
		fmt.Println("Commands: check, list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkCmd.Parse(os.Args[2:])
		if *checkURL == "" {
			fmt.Println("Error: --url is required")
			os.Exit(1)
		}
		checkService(*checkURL)
	case "list":
		listCmd.Parse(os.Args[2:])
		if *listFile == "" {
			fmt.Println("Error: --file is required")
			os.Exit(1)
		}
		listServices(*listFile)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func checkService(url string) {
	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start)

	if err != nil {
		fmt.Printf("FAIL: %s - %v\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	status := "OK"
	if resp.StatusCode >= 400 {
		status = "FAIL"
	}

	fmt.Printf("%s: %s (%d) - %s\n", url, status, resp.StatusCode, latency)
}

func listServices(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	_ = data
	fmt.Println("List services from file - implement parsing")
}
