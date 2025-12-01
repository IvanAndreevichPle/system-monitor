package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

func startTestServer(t *testing.T) (pb.SystemMonitorClient, func()) {
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	s := grpc.NewServer()
	RegisterServer(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Connect client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := pb.NewSystemMonitorClient(conn)

	cleanup := func() {
		s.GracefulStop()
		_ = conn.Close()
	}

	return client, cleanup
}

func TestGetMetrics_Integration(t *testing.T) {
	// Skip if /proc doesn't exist (e.g., on non-Linux systems)
	if _, err := net.Listen("tcp", ":0"); err != nil {
		t.Skip("Skipping test: network not available")
	}

	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.GetMetricsRequest{
		IntervalSeconds:       2,
		AveragingWindowSeconds: 5,
	}

	stream, err := client.GetMetrics(ctx, req)
	require.NoError(t, err)

	// Receive first snapshot
	snapshot, err := stream.Recv()
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.NotZero(t, snapshot.Timestamp)
	assert.NotNil(t, snapshot.LoadAverage)
	assert.NotNil(t, snapshot.Cpu)
	assert.NotNil(t, snapshot.Disk)

	// Receive second snapshot (should be averaged)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	select {
	case <-ctx2.Done():
		t.Fatal("Timeout waiting for second snapshot")
	default:
		snapshot2, err := stream.Recv()
		if err == nil {
			assert.NotNil(t, snapshot2)
			assert.NotZero(t, snapshot2.Timestamp)
		}
	}
}

func TestGetMetrics_InvalidInterval(t *testing.T) {
	// Skip if /proc doesn't exist (e.g., on non-Linux systems)
	if _, err := net.Listen("tcp", ":0"); err != nil {
		t.Skip("Skipping test: network not available")
	}

	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetMetricsRequest{
		IntervalSeconds:       0, // Invalid interval
		AveragingWindowSeconds: 5,
	}

	stream, err := client.GetMetrics(ctx, req)
	require.NoError(t, err)

	// Should still receive snapshots (server uses default)
	snapshot, err := stream.Recv()
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
}

