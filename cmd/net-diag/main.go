package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	checkTarget := checkCmd.String("target", "8.8.8.8", "Target IP or hostname to check")

	portsCmd := flag.NewFlagSet("ports", flag.ExitOnError)
	portsHost := portsCmd.String("host", "127.0.0.1", "Host to scan")

	dnsCmd := flag.NewFlagSet("dns", flag.ExitOnError)
	dnsDomain := dnsCmd.String("domain", "example.com", "Domain to resolve")

	if len(os.Args) < 2 {
		fmt.Println("Usage: net-diag <command> [arguments]")
		fmt.Println("Commands: check, ports, dns")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkCmd.Parse(os.Args[2:])
		fmt.Printf("=== Running Network Diagnostics for %s ===\n", *checkTarget)
		
		// TCP Ping test
		start := time.Now()
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(*checkTarget, "80"), 3*time.Second)
		duration := time.Since(start)
		if err != nil {
			// try port 443
			start = time.Now()
			conn, err = net.DialTimeout("tcp", net.JoinHostPort(*checkTarget, "443"), 3*time.Second)
			duration = time.Since(start)
		}

		if err != nil {
			fmt.Printf("[FAIL] TCP connectivity to %s failed: %v\n", *checkTarget, err)
		} else {
			conn.Close()
			fmt.Printf("[OK] TCP connectivity to %s succeeded in %v\n", *checkTarget, duration)
		}

		// DNS lookup check
		ips, err := net.LookupIP(*checkTarget)
		if err != nil {
			// If target is an IP, try reverse lookup or skip
			fmt.Printf("[INFO] DNS lookup not applicable or failed for %s: %v\n", *checkTarget, err)
		} else {
			fmt.Printf("[OK] Resolved IPs for %s: ", *checkTarget)
			for _, ip := range ips {
				fmt.Printf("%s ", ip.String())
			}
			fmt.Println()
		}

	case "ports":
		portsCmd.Parse(os.Args[2:])
		fmt.Printf("=== Scanning Common Homelab Ports on %s ===\n", *portsHost)
		commonPorts := []int{22, 53, 80, 443, 3000, 8006, 8080, 9090, 32400}

		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, port := range commonPorts {
			wg.Add(1)
			go func(p int) {
				defer wg.Done()
				address := fmt.Sprintf("%s:%d", *portsHost, p)
				conn, err := net.DialTimeout("tcp", address, 1*time.Second)
				if err == nil {
					conn.Close()
					mu.Lock()
					fmt.Printf("[OPEN] Port %d is open\n", p)
					mu.Unlock()
				}
			}(port)
		}
		wg.Wait()
		fmt.Println("Scan completed.")

	case "dns":
		dnsCmd.Parse(os.Args[2:])
		fmt.Printf("=== Performing DNS Lookup for %s ===\n", *dnsDomain)
		
		start := time.Now()
		records, err := net.LookupIP(*dnsDomain)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("[FAIL] Failed to resolve domain %s: %v\n", *dnsDomain, err)
			os.Exit(1)
		}

		fmt.Printf("[OK] Resolved in %v\n", duration)
		for _, r := range records {
			fmt.Printf("  -> IP: %s\n", r.String())
		}

		cnames, err := net.LookupCNAME(*dnsDomain)
		if err == nil && cnames != *dnsDomain {
			fmt.Printf("  -> CNAME: %s\n", cnames)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
