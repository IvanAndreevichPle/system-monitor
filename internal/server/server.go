package server

import (
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/IvanAndreevichPle/system-monitor/internal/collector"
	"github.com/IvanAndreevichPle/system-monitor/internal/storage"
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

	// Validate averaging window
	if averagingWindow <= 0 {
		averagingWindow = interval // default to interval if not specified
	}

	// Create storage for this client connection
	// Store metrics for at least averagingWindow + some buffer
	maxAge := averagingWindow
	if maxAge < 60 {
		maxAge = 60 // minimum 60 seconds
	}
	metricsStorage := storage.NewMetricsStorage(maxAge)

	// Create ticker for collecting metrics (collect more frequently for better averaging)
	// Collect every second or at interval, whichever is smaller
	collectInterval := interval
	if collectInterval > 1 {
		collectInterval = 1 // collect at least every second
	}
	collectTicker := time.NewTicker(time.Duration(collectInterval) * time.Second)
	defer collectTicker.Stop()

	// Create ticker for sending metrics
	sendTicker := time.NewTicker(time.Duration(interval) * time.Second)
	defer sendTicker.Stop()

	// Send initial snapshot immediately
	snapshot, err := s.collector.CollectAll()
	if err != nil {
		log.Printf("Failed to collect metrics: %v", err)
		return err
	}
	snapshot.Timestamp = time.Now().Unix()
	metricsStorage.AddSnapshot(snapshot)
	if err := stream.Send(snapshot); err != nil {
		log.Printf("Failed to send metrics: %v", err)
		return err
	}

	// Collect and send metrics
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client disconnected")
			return nil
		case <-collectTicker.C:
			// Collect metrics and store them
			snapshot, err := s.collector.CollectAll()
			if err != nil {
				log.Printf("Failed to collect metrics: %v", err)
				continue
			}
			snapshot.Timestamp = time.Now().Unix()
			metricsStorage.AddSnapshot(snapshot)
		case <-sendTicker.C:
			// Send averaged metrics
			averagedSnapshot := metricsStorage.GetAveragedSnapshot(averagingWindow)
			if averagedSnapshot == nil {
				log.Printf("No metrics available for averaging")
				continue
			}
			if err := stream.Send(averagedSnapshot); err != nil {
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

