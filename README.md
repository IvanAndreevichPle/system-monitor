# System Monitor

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

```bash
go build -o bin/monitor ./cmd/monitor
go build -o bin/client ./cmd/client
```

## Usage

### Start the daemon

```bash
./bin/monitor --port 8080
```

### Run the client

```bash
./bin/client --address localhost:8080
```

## Development

### Run tests

```bash
go test -race -count 100 ./...
```

### Run linter

```bash
golangci-lint run
```

## Project Structure

```
.
├── cmd/
│   ├── monitor/     # Main daemon
│   └── client/      # Simple client
├── internal/
│   ├── collector/   # Metric collectors
│   └── storage/      # In-memory metric storage
├── api/
│   └── proto/       # Protobuf definitions
└── .github/
    └── workflows/   # CI/CD pipeline
```

## License

MIT

