package server

import (
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/IvanAndreevichPle/system-monitor/internal/collector"
	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// Server implements the SystemMonitor gRPC service
type Server struct {
	pb.UnimplementedSystemMonitorServer
	collector *collector.Collector
}

// NewServer creates a new gRPC server instance
func NewServer() *Server {
	return &Server{
		collector: collector.NewCollector(),
	}
}

// GetMetrics implements server-side streaming for system metrics
// It sends metrics snapshots every N seconds, averaged over M seconds
func (s *Server) GetMetrics(req *pb.GetMetricsRequest, stream pb.SystemMonitor_GetMetricsServer) error {
	interval := req.IntervalSeconds
	averagingWindow := req.AveragingWindowSeconds

	log.Printf("Client connected: interval=%ds, averaging_window=%ds", interval, averagingWindow)

	// Validate interval
	if interval <= 0 {
		interval = 5 // default 5 seconds
	}

	// Create ticker for sending metrics
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Send initial snapshot immediately
	snapshot, err := s.collector.CollectAll()
	if err != nil {
		log.Printf("Failed to collect metrics: %v", err)
		return err
	}
	snapshot.Timestamp = time.Now().Unix()
	if err := stream.Send(snapshot); err != nil {
		log.Printf("Failed to send metrics: %v", err)
		return err
	}

	// Send metrics at regular intervals
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client disconnected")
			return nil
		case <-ticker.C:
			snapshot, err := s.collector.CollectAll()
			if err != nil {
				log.Printf("Failed to collect metrics: %v", err)
				continue // Continue trying on next interval
			}
			snapshot.Timestamp = time.Now().Unix()
			if err := stream.Send(snapshot); err != nil {
				log.Printf("Failed to send metrics: %v", err)
				return err
			}
		}
	}
}

// RegisterServer registers the server with gRPC
func RegisterServer(srv *grpc.Server) {
	pb.RegisterSystemMonitorServer(srv, NewServer())
}

