package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

func TestNewMetricsStorage(t *testing.T) {
	storage := NewMetricsStorage(60)
	assert.NotNil(t, storage)
}

func TestAddSnapshot(t *testing.T) {
	storage := NewMetricsStorage(60)

	snapshot := &pb.MetricsSnapshot{
		Timestamp: time.Now().Unix(),
		LoadAverage: &pb.LoadAverage{
			Load1:  1.0,
			Load5:  2.0,
			Load15: 3.0,
		},
	}

	storage.AddSnapshot(snapshot)
	// Should not panic and snapshot should be stored
}

func TestGetAveragedSnapshot_SingleSnapshot(t *testing.T) {
	storage := NewMetricsStorage(60)

	snapshot := &pb.MetricsSnapshot{
		Timestamp: time.Now().Unix(),
		LoadAverage: &pb.LoadAverage{
			Load1:  1.0,
			Load5:  2.0,
			Load15: 3.0,
		},
		Cpu: &pb.CPUMetrics{
			UserPercent:   10.0,
			SystemPercent: 20.0,
			IdlePercent:   70.0,
		},
	}

	storage.AddSnapshot(snapshot)
	averaged := storage.GetAveragedSnapshot(30)

	require.NotNil(t, averaged)
	assert.Equal(t, snapshot.LoadAverage.Load1, averaged.LoadAverage.Load1)
	assert.Equal(t, snapshot.Cpu.UserPercent, averaged.Cpu.UserPercent)
}

func TestGetAveragedSnapshot_MultipleSnapshots(t *testing.T) {
	storage := NewMetricsStorage(60)

	now := time.Now().Unix()

	// Add multiple snapshots
	for i := 0; i < 5; i++ {
		snapshot := &pb.MetricsSnapshot{
			Timestamp: now + int64(i),
			LoadAverage: &pb.LoadAverage{
				Load1:  float64(i + 1),
				Load5:   float64(i + 2),
				Load15: float64(i + 3),
			},
			Cpu: &pb.CPUMetrics{
				UserPercent:   float64(10 + i),
				SystemPercent: float64(20 + i),
				IdlePercent:   float64(70 - i),
			},
		}
		storage.AddSnapshot(snapshot)
	}

	averaged := storage.GetAveragedSnapshot(30)
	require.NotNil(t, averaged)

	// Average of 1, 2, 3, 4, 5 = 3.0
	expectedLoad1 := (1.0 + 2.0 + 3.0 + 4.0 + 5.0) / 5.0
	assert.InDelta(t, expectedLoad1, averaged.LoadAverage.Load1, 0.01)
}

func TestGetAveragedSnapshot_EmptyStorage(t *testing.T) {
	storage := NewMetricsStorage(60)
	averaged := storage.GetAveragedSnapshot(30)
	assert.Nil(t, averaged)
}

func TestCleanup(t *testing.T) {
	storage := NewMetricsStorage(5) // 5 seconds max age

	now := time.Now().Unix()

	// Add old snapshot
	oldSnapshot := &pb.MetricsSnapshot{
		Timestamp: now - 10, // 10 seconds ago
		LoadAverage: &pb.LoadAverage{
			Load1: 1.0,
		},
	}
	storage.AddSnapshot(oldSnapshot)

	// Add recent snapshot
	recentSnapshot := &pb.MetricsSnapshot{
		Timestamp: now,
		LoadAverage: &pb.LoadAverage{
			Load1: 2.0,
		},
	}
	storage.AddSnapshot(recentSnapshot)

	// Get averaged snapshot - should only include recent one
	averaged := storage.GetAveragedSnapshot(30)
	require.NotNil(t, averaged)
	assert.Equal(t, 2.0, averaged.LoadAverage.Load1)
}

