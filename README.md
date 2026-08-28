# Go CLI Tools

Command-line utilities written in Go for homelab management.

## Tools

| Tool | Description |
|------|-------------|
| `k8s-deploy` | Deploy/list/delete Kubernetes manifests |
| `docker-compose` | Docker Compose wrapper |
| `proxmox-cli` | Proxmox VE management |
| `net-diag` | Network diagnostics |
| `vault-cli` | HashiCorp Vault operations |
| `consul-cli` | Consul service management |
| `helm-cli` | Helm chart operations |
| `log-analyzer` | Log file analysis |
| `file-watcher` | File change monitoring |
| `json-formatter` | JSON pretty-printing |
| `process-monitor` | Process monitoring |

## Build

```bash
go build ./cmd/...
```

## Install

```bash
go install ./cmd/k8s-deploy
go install ./cmd/docker-compose
```

## Usage

```bash
k8s-deploy deploy --manifest app.yaml --namespace default
docker-compose up -d
log-analyzer --file app.log --pattern "ERROR"
```

## Requirements

- Go 1.21+

## License

MIT
