# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all tools
RUN mkdir -p /build && \
    for tool in proxmox-cli docker-compose net-diag k8s-deploy helm-cli vault-cli consul-cli; do \
        CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/${tool} ./cmd/${tool}; \
    done

# Runtime stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl bash

WORKDIR /usr/local/bin

# Copy all binaries
COPY --from=builder /build/* ./

# Create symlinks for convenience
RUN for tool in proxmox-cli docker-compose net-diag k8s-deploy helm-cli vault-cli consul-cli; do \
        ln -sf /usr/local/bin/${tool} /usr/local/bin/go-cli-${tool}; \
    done

ENTRYPOINT ["/usr/local/bin/proxmox-cli"]
LABEL org.opencontainers.image.source="https://github.com/magos-cyber/go-cli-tools"
LABEL org.opencontainers.image.description="Go CLI Tools for Homelab Automation"
LABEL org.opencontainers.image.licenses="MIT"