# System Monitor

[![CI](https://github.com/IvanAndreevichPle/system-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/IvanAndreevichPle/system-monitor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/IvanAndreevichPle/system-monitor)](https://goreportcard.com/report/github.com/IvanAndreevichPle/system-monitor)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Final project for OTUS Golang course: System monitoring daemon that collects system metrics and streams them to clients via gRPC.

## Description

A daemon that collects system information (load average, CPU, disk, network statistics) and sends it to clients using gRPC server-side streaming.

## Features

- **Load Average** - system load average from `/proc/loadavg`
- **CPU Metrics** - CPU usage (%user, %system, %idle) from `/proc/stat`
- **Disk Metrics** - disk I/O statistics (tps, KB/s) from `/proc/diskstats`
- **Filesystem Info** - disk usage and inode statistics
- **Network Statistics** - top talkers by protocol and traffic
- **Socket Statistics** - listening sockets and TCP connection states

## Requirements

- Go 1.22 or higher
- Linux (Ubuntu 18.04+)

## Installation

### Using Make

```bash
make build          # Build binaries
make install        # Install to /usr/local/bin
```

### Manual Build

```bash
go build -o bin/monitor ./cmd/monitor
go build -o bin/client ./cmd/client
```

### Using Docker

```bash
make docker-build   # Build Docker image
make docker-run     # Run in Docker container
```

## Usage

### Configuration

The daemon supports configuration via CLI flags and/or YAML config file. CLI flags override config file values.

#### Using CLI flags

```bash
./bin/monitor --port 8080 --log-level info --log-format text
```

#### Using config file

1. Copy example config:
```bash
cp configs/config.example.yaml configs/config.yaml
```

2. Edit `configs/config.yaml` as needed

3. Start with config file:
```bash
./bin/monitor --config configs/config.yaml
```

#### Available CLI flags

- `--port` - Port to listen on (default: 8080)
- `--log-level` - Log level: debug, info, warn, error (default: info)
- `--log-format` - Log format: json, text (default: text)
- `--config` - Path to config file

### Start the daemon

```bash
# Using defaults
./bin/monitor

# Using CLI flags
./bin/monitor --port 8080

# Using config file
./bin/monitor --config configs/config.yaml

# Mixing config file and CLI flags (CLI overrides config)
./bin/monitor --config configs/config.yaml --port 9090
```

### Run the client

```bash
./bin/client --address localhost:8080
```

## Development

### Using Make

```bash
make help           # Show all available commands
make build          # Build binaries
make test           # Run unit tests
make test-integration  # Run integration tests
make test-all       # Run all tests
make lint           # Run linter
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-run     # Run in Docker
```

### Manual Commands

#### Run tests

```bash
go test -race -count 100 ./...
go test -tags=integration ./internal/server/...
```

#### Run linter

```bash
golangci-lint run
```

#### Generate Protobuf code

```bash
make proto
# or manually:
protoc --go_out=. --go-grpc_out=. api/proto/metrics.proto
```

## Docker

The application requires access to `/proc` filesystem to collect system metrics. When running in Docker:

```bash
# Build image
make docker-build

# Run with host network and /proc mounted
make docker-run

# Or manually:
docker run -d \
  --name system-monitor \
  --network host \
  -v /proc:/host/proc:ro \
  -v /sys:/host/sys:ro \
  system-monitor:latest
```

**Note:** The application reads from `/proc` which is typically available on Linux hosts. For Docker, use `--network host` and mount `/proc` as read-only.

## Project Structure

```
.
├── cmd/
│   ├── monitor/     # Main daemon
│   └── client/      # Simple client
├── internal/
│   ├── collector/   # Metric collectors
│   ├── config/      # Configuration management
│   ├── server/      # gRPC server
│   └── storage/     # In-memory metric storage
├── api/
│   └── proto/       # Protobuf definitions
├── configs/
│   └── config.example.yaml  # Example configuration file
├── Dockerfile       # Docker image definition
├── Makefile        # Build automation
└── .github/
    └── workflows/   # CI/CD pipeline
```

## License

MIT

