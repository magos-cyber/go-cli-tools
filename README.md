# Go CLI Tools for Homelab Automation

Production-ready Go CLI tools for managing homelab infrastructure — Proxmox VE, Docker, Kubernetes, Helm, Vault, Consul, and network diagnostics.

## 📁 Structure

```
go-cli-tools/
├── cmd/
│   ├── proxmox-cli/       # Proxmox VE management CLI
│   ├── docker-compose/    # Docker Compose stack manager
│   ├── net-diag/          # Network diagnostics tool
│   ├── k8s-deploy/        # Kubernetes deployment manager
│   ├── helm-cli/          # Helm chart manager
│   ├── vault-cli/         # HashiCorp Vault client
│   └── consul-cli/        # Consul service discovery client
├── .github/
│   └── workflows/
│       └── release.yml    # CI/CD pipeline
├── go.mod
├── go.sum
├── Makefile
└── LICENSE
```

## 🚀 Quick Start

```bash
# Clone the repo
git clone https://github.com/magos-cyber/go-cli-tools.git
cd go-cli-tools

# Build all tools
make build

# Or build individual tools
make proxmox-cli
make docker-compose
make net-diag
make k8s-deploy
make helm-cli
make vault-cli
make consul-cli

# Install to $GOPATH/bin
make install
```

## 📝 Tools

### proxmox-cli
```bash
# List all VMs and containers
./proxmox-cli list --node pve

# Start a VM
./proxmox-cli start --node pve --vmid 100 --type qemu

# Stop a VM
./proxmox-cli stop --node pve --vmid 100 --type qemu

# Backup a VM
./proxmox-cli backup --node pve --vmid 100 --storage local

# Cluster status
./proxmox-cli status
```

### docker-compose
```bash
# Deploy a stack
./docker-compose deploy --stack vaultwarden

# List active stacks
./docker-compose list

# Update a stack
./docker-compose update --stack vaultwarden

# Update all stacks
./docker-compose update --all

# View logs
./docker-compose logs --stack vaultwarden --tail 50
```

### net-diag
```bash
# Check connectivity
./net-diag check --target 8.8.8.8

# Scan common ports
./net-diag ports --host 10.0.0.10

# DNS lookup
./net-diag dns --domain example.com
```

### k8s-deploy
```bash
# Deploy a manifest
./k8s-deploy deploy --manifest ./deployment.yaml --namespace default

# List resources
./k8s-deploy list --namespace default

# Delete a resource
./k8s-deploy delete --resource deployment/nginx --namespace default
```

### helm-cli
```bash
# Install a chart
./helm-cli install --name my-release --chart nginx-ingress --namespace default

# List releases
./helm-cli list

# Upgrade a release
./helm-cli upgrade --name my-release --chart nginx-ingress --namespace default

# Uninstall a release
./helm-cli uninstall --name my-release --namespace default
```

### vault-cli
```bash
# Read a secret
./vault-cli read --path secret/data/myapp

# Write a secret
./vault-cli write --path secret/data/myapp --data "key1=value1,key2=value2"

# List secrets
./vault-cli list --path secret/

# Delete a secret
./vault-cli delete --path secret/data/myapp
```

### consul-cli
```bash
# Register a service
./consul-cli register --name web --id web-1 --address 10.0.0.50 --port 8080

# Deregister a service
./consul-cli deregister --id web-1

# List services
./consul-cli list
```

## 🔧 Configuration

Most tools use environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PVE_HOST` | Proxmox API endpoint | `https://localhost:8006/api2/json` |
| `PVE_TOKEN_ID` | Proxmox API token ID | - |
| `PVE_SECRET` | Proxmox API secret | - |
| `DOCKER_STACKS_DIR` | Docker stacks directory | `/opt/stacks` |

## 📦 Releases

Pre-built binaries for Linux, macOS, and Windows are available on the [Releases](https://github.com/magos-cyber/go-cli-tools/releases) page.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.