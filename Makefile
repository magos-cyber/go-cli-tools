.PHONY: all build test vet clean install

TOOLS := proxmox-cli docker-compose net-diag k8s-deploy helm-cli vault-cli consul-cli

all: build

build: $(TOOLS)

$(TOOLS):
	@echo "Building $@..."
	go build -o bin/$@ ./cmd/$@

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin/ release/

install: build
	@echo "Installing tools to $(GOPATH)/bin..."
	for tool in $(TOOLS); do \
		cp bin/$$tool $(GOPATH)/bin/; \
	done
	@echo "All tools installed to $(GOPATH)/bin"

.DEFAULT_GOAL := build