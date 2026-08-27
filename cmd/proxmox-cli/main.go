package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL, tokenId, secret string, verifySSL bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL},
	}
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Token:   fmt.Sprintf("PVEAPIToken=%s=%s", tokenId, secret),
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		},
	}
}

func (c *Client) DoRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	reqURL := c.BaseURL + endpoint
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %string", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

type Response struct {
	Data json.RawMessage `json:"data"`
}

func main() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listNode := listCmd.String("node", "pve", "Proxmox node name")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startNode := startCmd.String("node", "pve", "Proxmox node name")
	startVMID := startCmd.Int("vmid", 0, "VM or Container ID")
	startType := startCmd.String("type", "qemu", "Type: qemu or lxc")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopNode := stopCmd.String("node", "pve", "Proxmox node name")
	stopVMID := stopCmd.Int("vmid", 0, "VM or Container ID")
	stopType := stopCmd.String("type", "qemu", "Type: qemu or lxc")

	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	backupNode := backupCmd.String("node", "pve", "Proxmox node name")
	backupVMID := backupCmd.Int("vmid", 0, "VM or Container ID")
	backupStorage := backupCmd.String("storage", "local", "Storage target for backup")

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: proxmox-cli <command> [arguments]")
		fmt.Println("Commands: list, start, stop, backup, status")
		os.Exit(1)
	}

	endpoint := os.Getenv("PVE_HOST")
	if endpoint == "" {
		endpoint = "https://localhost:8006/api2/json"
	}
	tokenId := os.Getenv("PVE_TOKEN_ID")
	secret := os.Getenv("PVE_SECRET")
	verifySSL := os.Getenv("PVE_VERIFY_SSL") == "true"

	client := NewClient(endpoint, tokenId, secret, verifySSL)

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		// Fetch QEMU VMs and LXC Containers
		qemuBytes, err := client.DoRequest("GET", fmt.Sprintf("/nodes/%s/qemu", *listNode), nil)
		if err != nil {
			fmt.Printf("Error fetching VMs: %v\n", err)
		}
		lxcBytes, err := client.DoRequest("GET", fmt.Sprintf("/nodes/%s/lxc", *listNode), nil)
		if err != nil {
			fmt.Printf("Error fetching LXC: %v\n", err)
		}

		fmt.Printf("=== Proxmox Resources on Node: %s ===\n", *listNode)
		printVMs("QEMU VMs", qemuBytes)
		printVMs("LXC Containers", lxcBytes)

	case "start":
		startCmd.Parse(os.Args[2:])
		if *startVMID == 0 {
			fmt.Println("Error: --vmid is required")
			os.Exit(1)
		}
		path := fmt.Sprintf("/nodes/%s/%s/%d/status/start", *startNode, *startType, *startVMID)
		_, err := client.DoRequest("POST", path, nil)
		if err != nil {
			fmt.Printf("Failed to start VM/LXC %d: %v\n", *startVMID, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully sent start command to %s %d\n", *startType, *startVMID)

	case "stop":
		stopCmd.Parse(os.Args[2:])
		if *stopVMID == 0 {
			fmt.Println("Error: --vmid is required")
			os.Exit(1)
		}
		path := fmt.Sprintf("/nodes/%s/%s/%d/status/stop", *stopNode, *stopType, *stopVMID)
		_, err := client.DoRequest("POST", path, nil)
		if err != nil {
			fmt.Printf("Failed to stop VM/LXC %d: %v\n", *stopVMID, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully sent stop command to %s %d\n", *stopType, *stopVMID)

	case "backup":
		backupCmd.Parse(os.Args[2:])
		if *backupVMID == 0 {
			fmt.Println("Error: --vmid is required")
			os.Exit(1)
		}
		path := fmt.Sprintf("/nodes/%s/vzdump", *backupNode)
		data := url.Values{}
		data.Set("vmid", fmt.Sprintf("%d", *backupVMID))
		data.Set("storage", *backupStorage)
		data.Set("mode", "snapshot")

		_, err := client.DoRequest("POST", path, strings.NewReader(data.Encode()))
		if err != nil {
			fmt.Printf("Failed to initiate backup for VM %d: %v\n", *backupVMID, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully initiated backup for VM %d to storage %s\n", *backupVMID, *backupStorage)

	case "status":
		statusCmd.Parse(os.Args[2:])
		clusterBytes, err := client.DoRequest("GET", "/cluster/status", nil)
		if err != nil {
			fmt.Printf("Failed to fetch cluster status: %v\n", err)
			os.Exit(1)
		}
		var resp Response
		if err := json.Unmarshal(clusterBytes, &resp); err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
			os.Exit(1)
		}
		var nodes []map[string]interface{}
		json.Unmarshal(resp.Data, &nodes)
		fmt.Println("=== Proxmox Cluster Status ===")
		for _, n := range nodes {
			fmt.Printf("- Name: %v | Type: %v | Status: %v\n", n["name"], n["type"], n["status"])
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printVMs(title string, raw []byte) {
	var resp Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		fmt.Printf("Could not parse %s\n", title)
		return
	}
	var items []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &items); err != nil {
		return
	}
	fmt.Printf("\n--- %s ---\n", title)
	for _, item := range items {
		fmt.Printf("ID: %v | Name: %v | Status: %v | CPU: %v | Mem: %v\n",
			item["vmid"], item["name"], item["status"], item["cpu"], item["mem"])
	}
}
