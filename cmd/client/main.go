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
	defer conn.Close()

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

		// TODO: Print actual metrics when implemented
		if snapshot.LoadAverage != nil {
			fmt.Printf("Load Average: %.2f, %.2f, %.2f\n",
				snapshot.LoadAverage.Load1,
				snapshot.LoadAverage.Load5,
				snapshot.LoadAverage.Load15)
		}

		if snapshot.Cpu != nil {
			fmt.Printf("CPU: user=%.2f%%, system=%.2f%%, idle=%.2f%%\n",
				snapshot.Cpu.UserPercent,
				snapshot.Cpu.SystemPercent,
				snapshot.Cpu.IdlePercent)
		}
	}
}

