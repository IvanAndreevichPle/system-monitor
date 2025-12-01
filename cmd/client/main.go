package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

var (
	address = flag.String("address", "localhost:8080", "Server address")
)

func main() {
	flag.Parse()

	// Connect to server
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	client := pb.NewSystemMonitorClient(conn)

	// Create request
	req := &pb.GetMetricsRequest{
		IntervalSeconds:       5,
		AveragingWindowSeconds: 15,
	}

	// Create context with cancellation
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Get metrics stream
	stream, err := client.GetMetrics(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get metrics: %v", err)
	}

	log.Printf("Connected to server at %s", *address)
	log.Printf("Receiving metrics (interval=%ds, window=%ds)...", req.IntervalSeconds, req.AveragingWindowSeconds)

	// Receive metrics
	for {
		snapshot, err := stream.Recv()
		if err != nil {
			log.Printf("Stream closed: %v", err)
			return
		}

		// Print timestamp
		timestamp := time.Unix(snapshot.Timestamp, 0)
		fmt.Printf("\n=== Metrics Snapshot at %s ===\n", timestamp.Format(time.RFC3339))

		// Print Load Average
		if snapshot.LoadAverage != nil {
			fmt.Printf("Load Average: %.2f (1m), %.2f (5m), %.2f (15m)\n",
				snapshot.LoadAverage.Load1,
				snapshot.LoadAverage.Load5,
				snapshot.LoadAverage.Load15)
		}

		// Print CPU metrics
		if snapshot.Cpu != nil {
			fmt.Printf("CPU: user=%.2f%%, system=%.2f%%, idle=%.2f%%\n",
				snapshot.Cpu.UserPercent,
				snapshot.Cpu.SystemPercent,
				snapshot.Cpu.IdlePercent)
		}

		// Print Disk metrics
		if snapshot.Disk != nil {
			fmt.Printf("Disk: TPS=%.2f, KB/s=%.2f\n", snapshot.Disk.Tps, snapshot.Disk.KbPerSec)
			if len(snapshot.Disk.Disks) > 0 {
				fmt.Printf("  Per-disk stats:\n")
				for _, disk := range snapshot.Disk.Disks {
					fmt.Printf("    %s: TPS=%.2f, KB/s=%.2f\n", disk.Device, disk.Tps, disk.KbPerSec)
				}
			}
		}

		// Print Filesystem info
		if snapshot.Filesystem != nil && len(snapshot.Filesystem.Filesystems) > 0 {
			fmt.Printf("Filesystems:\n")
			for _, fs := range snapshot.Filesystem.Filesystems {
				fmt.Printf("  %s: %.2f%% used (%d MB), %.2f%% inodes used (%d)\n",
					fs.MountPoint, fs.UsedPercent, fs.UsedMb, fs.InodesPercent, fs.UsedInodes)
			}
		}

		// Print Network stats
		if snapshot.Network != nil {
			if len(snapshot.Network.ProtocolStats) > 0 {
				fmt.Printf("Network protocols:\n")
				for _, proto := range snapshot.Network.ProtocolStats {
					fmt.Printf("  %s: %d bytes (%.2f%%)\n", proto.Protocol, proto.Bytes, proto.Percent)
				}
			}
		}

		// Print Socket stats
		if snapshot.Sockets != nil {
			if snapshot.Sockets.TcpStates != nil {
				fmt.Printf("TCP States: ESTABLISHED=%d, TIME_WAIT=%d, FIN_WAIT=%d, SYN_RECEIVED=%d, CLOSE_WAIT=%d, OTHER=%d\n",
					snapshot.Sockets.TcpStates.Established,
					snapshot.Sockets.TcpStates.TimeWait,
					snapshot.Sockets.TcpStates.FinWait,
					snapshot.Sockets.TcpStates.SynReceived,
					snapshot.Sockets.TcpStates.CloseWait,
					snapshot.Sockets.TcpStates.Other)
			}
			if len(snapshot.Sockets.ListeningSockets) > 0 {
				fmt.Printf("Listening sockets: %d\n", len(snapshot.Sockets.ListeningSockets))
			}
		}
	}
}

