package server

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// Server implements the SystemMonitor gRPC service
type Server struct {
	pb.UnimplementedSystemMonitorServer
}

// NewServer creates a new gRPC server instance
func NewServer() *Server {
	return &Server{}
}

// GetMetrics implements server-side streaming for system metrics
// It sends metrics snapshots every N seconds, averaged over M seconds
func (s *Server) GetMetrics(req *pb.GetMetricsRequest, stream pb.SystemMonitor_GetMetricsServer) error {
	interval := req.IntervalSeconds
	averagingWindow := req.AveragingWindowSeconds

	log.Printf("Client connected: interval=%ds, averaging_window=%ds", interval, averagingWindow)

	// TODO: Implement metrics collection and averaging
	// For now, send empty snapshots as a placeholder
	ticker := stream.Context().Done()
	
	// Wait for context cancellation (client disconnect)
	<-ticker
	log.Printf("Client disconnected")
	
	return nil
}

// RegisterServer registers the server with gRPC
func RegisterServer(srv *grpc.Server) {
	pb.RegisterSystemMonitorServer(srv, NewServer())
}

