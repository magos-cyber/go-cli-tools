# Go CLI Tools for Homelab Automation

Production-ready Go CLI tools for managing homelab infrastructure — Proxmox VE, Docker, Kubernetes, Helm, Vault, Consul, and network diagnostics.

## [FOLDER] Structure

```
go-cli-tools/
+-- cmd/
|   +-- proxmox-cli/       # Proxmox VE management CLI
|   +-- docker-compose/    # Docker Compose stack manager
|   +-- net-diag/          # Network diagnostics tool
|   +-- k8s-deploy/        # Kubernetes deployment manager
|   +-- helm-cli/          # Helm chart manager
|   +-- vault-cli/         # HashiCorp Vault client
|   `-- consul-cli/        # Consul service discovery client
+-- .github/
|   `-- workflows/
|       +-- release.yml    # Build & release binaries
|       `-- docker.yml     # Build & push Docker image
+-- go.mod
+-- go.sum
+-- Dockerfile
+-- Makefile
`-- LICENSE
```

## [ROCKET] Quick Start

### Pre-built binaries
```bash
# Download latest release
curl -sL https://github.com/magos-cyber/go-cli-tools/releases/latest/download/go-cli-tools-v1.0.0-linux-amd64.tar.gz | tar xz
```

### Docker
```bash
docker run -it --rm ghcr.io/magos-cyber/go-cli-tools proxmox-cli list --node pve
```

### Build from source
```bash
git clone https://github.com/magos-cyber/go-cli-tools.git
cd go-cli-tools
make build
```

## [MEMO] Tools

### proxmox-cli
```bash
./proxmox-cli list --node pve
./proxmox-cli start --node pve --vmid 100 --type qemu
./proxmox-cli stop --node pve --vmid 100 --type qemu
./proxmox-cli backup --node pve --vmid 100 --storage local
./proxmox-cli status
```

### docker-compose
```bash
./docker-compose deploy --stack vaultwarden
./docker-compose list
./docker-compose update --stack vaultwarden
./docker-compose update --all
./docker-compose logs --stack vaultwarden --tail 50
```

### net-diag
```bash
./net-diag check --target 8.8.8.8
./net-diag ports --host 10.0.0.10
./net-diag dns --domain example.com
```

### k8s-deploy
```bash
./k8s-deploy deploy --manifest ./deployment.yaml --namespace default
./k8s-deploy list --namespace default
./k8s-deploy delete --resource deployment/nginx --namespace default
```

### helm-cli
```bash
./helm-cli install --name my-release --chart nginx-ingress --namespace default
./helm-cli list
./helm-cli upgrade --name my-release --chart nginx-ingress --namespace default
./helm-cli uninstall --name my-release --namespace default
```

### vault-cli
```bash
./vault-cli read --path secret/data/myapp
./vault-cli write --path secret/data/myapp --data "key1=value1,key2=value2"
./vault-cli list --path secret/
./vault-cli delete --path secret/data/myapp
```

### consul-cli
```bash
./consul-cli register --name web --id web-1 --address 10.0.0.50 --port 8080
./consul-cli deregister --id web-1
./consul-cli list
```

## [WRENCH] Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PVE_HOST` | Proxmox API endpoint | `https://localhost:8006/api2/json` |
| `PVE_TOKEN_ID` | Proxmox API token ID | - |
| `PVE_SECRET` | Proxmox API secret | - |
| `DOCKER_STACKS_DIR` | Docker stacks directory | `/opt/stacks` |

## [PACKAGE] Packages

Docker images are published to GitHub Container Registry:

| Package | Pull Command |
|---------|--------------|
| go-cli-tools | `docker pull ghcr.io/magos-cyber/go-cli-tools:latest` |

Tags:
- `:latest` — latest main branch build
- `:v1.0.0` — versioned releases
- `:sha-abc1234` — commit-specific builds

## [PAGE] License

MIT License - see [LICENSE](LICENSE) for details.