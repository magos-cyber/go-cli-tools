package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	cpuCmd := flag.NewFlagSet("cpu", flag.ExitOnError)
	memCmd := flag.NewFlagSet("mem", flag.ExitOnError)
	diskCmd := flag.NewFlagSet("disk", flag.ExitOnError)
	netCmd := flag.NewFlagSet("net", flag.ExitOnError)
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: system-monitor <command> [arguments]")
		fmt.Println("Commands: cpu, mem, disk, net, all")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "cpu":
		cpuCmd.Parse(os.Args[2:])
		showCPU()
	case "mem":
		memCmd.Parse(os.Args[2:])
		showMemory()
	case "disk":
		diskCmd.Parse(os.Args[2:])
		showDisk()
	case "net":
		netCmd.Parse(os.Args[2:])
		showNetwork()
	case "all":
		allCmd.Parse(os.Args[2:])
		showAll()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func showCPU() {
	fmt.Println("=== CPU Information ===")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("sh", "-c", "cat /proc/loadavg").Output()
		fields := strings.Fields(string(out))
		if len(fields) >= 3 {
			fmt.Printf("Load Average: %s, %s, %s\n", fields[0], fields[1], fields[2])
		}
	}
}

func showMemory() {
	fmt.Println("=== Memory Information ===")
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("free", "-h").Output()
		fmt.Println(string(out))
	}
}

func showDisk() {
	fmt.Println("=== Disk Usage ===")
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("df", "-h").Output()
		fmt.Println(string(out))
	}
}

func showNetwork() {
	fmt.Println("=== Network Interfaces ===")
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("ip", "-4", "addr", "show").Output()
		fmt.Println(string(out))
	}
}

func showAll() {
	fmt.Printf("System Monitor - %s\n\n", time.Now().Format(time.RFC3339))
	showCPU()
	fmt.Println()
	showMemory()
	fmt.Println()
	showDisk()
	fmt.Println()
	showNetwork()
}
