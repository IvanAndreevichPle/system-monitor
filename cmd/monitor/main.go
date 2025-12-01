package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/IvanAndreevichPle/system-monitor/internal/config"
	"github.com/IvanAndreevichPle/system-monitor/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register service
	server.RegisterServer(s)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-sigChan
		log.Println("Shutting down server...")
		s.GracefulStop()
	}()

	log.Printf("Server listening on port %d", cfg.Server.Port)
	log.Printf("Log level: %s, format: %s", cfg.Logging.Level, cfg.Logging.Format)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

