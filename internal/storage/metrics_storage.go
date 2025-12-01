package storage

import (
	"sync"
	"time"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// MetricsStorage stores metrics snapshots for averaging
type MetricsStorage struct {
	mu       sync.RWMutex
	snapshots []*pb.MetricsSnapshot
	maxAge   time.Duration
}

// NewMetricsStorage creates a new metrics storage with specified max age
func NewMetricsStorage(maxAgeSeconds int32) *MetricsStorage {
	return &MetricsStorage{
		snapshots: make([]*pb.MetricsSnapshot, 0),
		maxAge:    time.Duration(maxAgeSeconds) * time.Second,
	}
}

// AddSnapshot adds a new metrics snapshot to storage
func (s *MetricsStorage) AddSnapshot(snapshot *pb.MetricsSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add timestamp if not set
	if snapshot.Timestamp == 0 {
		snapshot.Timestamp = time.Now().Unix()
	}

	// Add snapshot
	s.snapshots = append(s.snapshots, snapshot)

	// Clean up old snapshots
	s.cleanup()
}

// cleanup removes snapshots older than maxAge
func (s *MetricsStorage) cleanup() {
	now := time.Now().Unix()
	cutoff := now - int64(s.maxAge.Seconds())

	// Find first snapshot within the window
	startIdx := 0
	for i, snap := range s.snapshots {
		if snap.Timestamp >= cutoff {
			startIdx = i
			break
		}
	}

	// Keep only snapshots within the window
	if startIdx > 0 {
		s.snapshots = s.snapshots[startIdx:]
	}
}

// GetAveragedSnapshot returns an averaged snapshot over the specified window
func (s *MetricsStorage) GetAveragedSnapshot(windowSeconds int32) *pb.MetricsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.snapshots) == 0 {
		return nil
	}

	now := time.Now().Unix()
	cutoff := now - int64(windowSeconds)

	// Filter snapshots within the window
	var validSnapshots []*pb.MetricsSnapshot
	for _, snap := range s.snapshots {
		if snap.Timestamp >= cutoff {
			validSnapshots = append(validSnapshots, snap)
		}
	}

	if len(validSnapshots) == 0 {
		// Return the most recent snapshot if no snapshots in window
		return s.snapshots[len(s.snapshots)-1]
	}

	// Average all metrics
	return averageSnapshots(validSnapshots)
}

// averageSnapshots calculates average values across multiple snapshots
func averageSnapshots(snapshots []*pb.MetricsSnapshot) *pb.MetricsSnapshot {
	if len(snapshots) == 0 {
		return nil
	}

	if len(snapshots) == 1 {
		return snapshots[0]
	}

	avg := &pb.MetricsSnapshot{
		Timestamp: time.Now().Unix(),
	}

	// Average Load Average
	if len(snapshots) > 0 && snapshots[0].LoadAverage != nil {
		var load1, load5, load15 float64
		count := 0
		for _, snap := range snapshots {
			if snap.LoadAverage != nil {
				load1 += snap.LoadAverage.Load1
				load5 += snap.LoadAverage.Load5
				load15 += snap.LoadAverage.Load15
				count++
			}
		}
		if count > 0 {
			avg.LoadAverage = &pb.LoadAverage{
				Load1:  load1 / float64(count),
				Load5:  load5 / float64(count),
				Load15: load15 / float64(count),
			}
		}
	}

	// Average CPU metrics
	if len(snapshots) > 0 && snapshots[0].Cpu != nil {
		var user, system, idle float64
		count := 0
		for _, snap := range snapshots {
			if snap.Cpu != nil {
				user += snap.Cpu.UserPercent
				system += snap.Cpu.SystemPercent
				idle += snap.Cpu.IdlePercent
				count++
			}
		}
		if count > 0 {
			avg.Cpu = &pb.CPUMetrics{
				UserPercent:   user / float64(count),
				SystemPercent: system / float64(count),
				IdlePercent:   idle / float64(count),
			}
		}
	}

	// Average Disk metrics
	if len(snapshots) > 0 && snapshots[0].Disk != nil {
		var totalTPS, totalKBPerSec float64
		diskMap := make(map[string]*diskAccumulator)
		count := 0

		for _, snap := range snapshots {
			if snap.Disk != nil {
				totalTPS += snap.Disk.Tps
				totalKBPerSec += snap.Disk.KbPerSec
				count++

				// Accumulate per-disk stats
				for _, disk := range snap.Disk.Disks {
					if acc, exists := diskMap[disk.Device]; exists {
						acc.tps += disk.Tps
						acc.kbPerSec += disk.KbPerSec
						acc.count++
					} else {
						diskMap[disk.Device] = &diskAccumulator{
							device:    disk.Device,
							tps:       disk.Tps,
							kbPerSec:  disk.KbPerSec,
							count:     1,
						}
					}
				}
			}
		}

		if count > 0 {
			disks := make([]*pb.DiskStats, 0, len(diskMap))
			for _, acc := range diskMap {
				disks = append(disks, &pb.DiskStats{
					Device:   acc.device,
					Tps:      acc.tps / float64(acc.count),
					KbPerSec: acc.kbPerSec / float64(acc.count),
				})
			}
			avg.Disk = &pb.DiskMetrics{
				Tps:      totalTPS / float64(count),
				KbPerSec: totalKBPerSec / float64(count),
				Disks:    disks,
			}
		}
	}

	// For filesystem, network, and sockets, use the most recent snapshot
	// as averaging these is more complex and may not be meaningful
	if len(snapshots) > 0 {
		latest := snapshots[len(snapshots)-1]
		avg.Filesystem = latest.Filesystem
		avg.Network = latest.Network
		avg.Sockets = latest.Sockets
	}

	return avg
}

// diskAccumulator accumulates disk statistics for averaging
type diskAccumulator struct {
	device    string
	tps       float64
	kbPerSec  float64
	count     int
}

