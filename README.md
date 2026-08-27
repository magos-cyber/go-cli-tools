# Go CLI Tools for Homelab Automation

Production-ready Go CLI tools for managing homelab infrastructure — Proxmox VE, Docker, network diagnostics, and monitoring.

## 📁 Structure

```
go-cli-tools/
├── cmd/
│   ├── proxmox-cli/       # Proxmox VE management CLI
│   ├── docker-compose/    # Docker Compose stack manager
│   ├── net-diag/          # Network diagnostics tool
│   └── sys-mon/           # System monitoring agent
├── internal/
│   ├── api/               # API client libraries
│   ├── config/            # Configuration management
│   └── utils/             # Shared utilities
├── scripts/
│   └── install.sh         # Build and install script
├── go.mod
├── go.sum
└── Makefile
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
make sys-mon
```

## 📝 Tools

### proxmox-cli
```bash
# List all VMs and containers
./proxmox-cli list --node pve

# Start a VM
./proxmox-cli vm start --node pve --vmid 100

# Create backup
./proxmox-cli backup --node pve --vmid 100 --storage backup-nfs

# Show cluster status
./proxmox-cli cluster status
```

### docker-compose
```bash
# Deploy a stack
./docker-compose deploy --stack monitoring --env-file .env

# List running stacks
./docker-compose list

# Update all stacks
./docker-compose update --all

# Show logs
./docker-compose logs --stack monitoring --tail 100
```

### net-diag
```bash
# Full network diagnostics
./net-diag check --target 8.8.8.8

# Port scan common homelab ports
./net-diag ports --host 192.168.1.1

# DNS resolution test
./net-diag dns --domain example.com

# Bandwidth test
./net-diag bandwidth --server speedtest.waust.org
```

### sys-mon
```bash
# System overview
./sys-mon overview

# Continuous monitoring
./sys-mon watch --interval 5s

# Export metrics to Prometheus
./sys-mon prometheus --port 9090

# Alert on thresholds
./sys-mon alert --cpu 90 --mem 85 --disk 95
```

## ⚙️ Configuration

```yaml
# ~/.config/go-cli-tools/config.yaml
proxmox:
  host: pve.example.com
  user: root@pam
  token: "your-api-token"
  verify_ssl: false

docker:
  socket: /var/run/docker.sock
  stacks_dir: /opt/stacks

monitoring:
  interval: 30s
  telegram:
    bot_token: "your-bot-token"
    chat_id: "your-chat-id"

prometheus:
  enabled: true
  port: 9090
```

## 🤝 Contributing

Contributions are welcome! Please:
- Write clean, idiomatic Go
- Add tests for new features
- Update documentation

## 📜 License

MIT License — see [LICENSE](LICENSE) for details.