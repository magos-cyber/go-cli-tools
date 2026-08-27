.PHONY: all build proxmox-cli docker-compose net-diag clean

all: build

build: proxmox-cli docker-compose net-diag

proxmox-cli:
	@echo "Building proxmox-cli..."
	go build -o bin/proxmox-cli ./cmd/proxmox-cli

docker-compose:
	@echo "Building docker-compose..."
	go build -o bin/docker-compose ./cmd/docker-compose

net-diag:
	@echo "Building net-diag..."
	go build -o bin/net-diag ./cmd/net-diag

clean:
	@echo "Cleaning binaries..."
	rm -rf bin/
