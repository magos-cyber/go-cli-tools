package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type CertInfo struct {
	Hostname  string
	Issuer    string
	Expires   string
	DaysLeft  int
	Status    string
}

func main() {
	hostCmd := flag.NewFlagSet("check", flag.ExitOnError)
	hostHostname := hostCmd.String("host", "", "Hostname to check")
	hostPort := hostCmd.Int("port", 443, "Port to check")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listFile := listCmd.String("file", "", "File with hosts")

	if len(os.Args) < 2 {
		fmt.Println("Usage: cert-manager <command> [arguments]")
		fmt.Println("Commands: check, list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		hostCmd.Parse(os.Args[2:])
		if *hostHostname == "" {
			fmt.Println("Error: --host is required")
			os.Exit(1)
		}
		doCheckCert(*hostHostname, *hostPort)
	case "list":
		listCmd.Parse(os.Args[2:])
		if *listFile == "" {
			fmt.Println("Error: --file is required")
			os.Exit(1)
		}
		listCerts(*listFile)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func doCheckCert(hostname string, port int) {
	addr := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		fmt.Println("No certificates found")
		os.Exit(1)
	}

	cert := certs[0]
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
	status := "OK"
	if daysLeft <= 0 {
		status = "EXPIRED"
	} else if daysLeft <= 7 {
		status = "CRITICAL"
	} else if daysLeft <= 30 {
		status = "WARNING"
	}

	info := CertInfo{
		Hostname: hostname,
		Issuer:   cert.Issuer.CommonName,
		Expires:  cert.NotAfter.Format("2006-01-02"),
		DaysLeft: daysLeft,
		Status:   status,
	}

	data, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(data))
}

func listCerts(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	_ = data
	fmt.Println("List certs from file - implement parsing")
}
