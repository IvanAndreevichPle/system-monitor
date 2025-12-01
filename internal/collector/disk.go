package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// diskStats holds raw disk statistics from /proc/diskstats
type diskStats struct {
	device         string
	readsCompleted uint64
	writesCompleted uint64
	sectorsRead     uint64
	sectorsWritten  uint64
}

// parseDiskStats parses a line from /proc/diskstats
// Format: "major minor device reads_completed reads_merged sectors_read time_reading writes_completed writes_merged sectors_written time_writing ..."
func parseDiskStats(line string) (*diskStats, error) {
	fields := strings.Fields(line)
	if len(fields) < 14 {
		return nil, fmt.Errorf("invalid /proc/diskstats format: expected at least 14 fields, got %d", len(fields))
	}

	device := fields[2]

	readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reads_completed: %w", err)
	}

	sectorsRead, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sectors_read: %w", err)
	}

	writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse writes_completed: %w", err)
	}

	sectorsWritten, err := strconv.ParseUint(fields[9], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sectors_written: %w", err)
	}

	return &diskStats{
		device:          device,
		readsCompleted:  readsCompleted,
		writesCompleted: writesCompleted,
		sectorsRead:     sectorsRead,
		sectorsWritten:  sectorsWritten,
	}, nil
}

// CollectDiskMetrics reads and parses /proc/diskstats
// Returns disk metrics with TPS and KB/s per disk
// Note: For accurate TPS and KB/s, we need two samples. This implementation
// returns current values. Averaging will be implemented in stage 10.
func (c *Collector) CollectDiskMetrics() (*pb.DiskMetrics, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/diskstats: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var totalTPS float64
	var totalKBPerSec float64
	var diskStatsList []*pb.DiskStats

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stats, err := parseDiskStats(line)
		if err != nil {
			// Skip invalid lines
			continue
		}

		// Skip loop devices and other virtual devices for now
		// We can filter them out or include them based on requirements
		if strings.HasPrefix(stats.device, "loop") {
			continue
		}

		// Calculate TPS (transfers per second) - for now, just use current values
		// In stage 10, we'll calculate the difference between two samples
		tps := float64(stats.readsCompleted + stats.writesCompleted)

		// Calculate KB/s (kilobytes per second)
		// Each sector is typically 512 bytes = 0.5 KB
		kbPerSec := float64(stats.sectorsRead+stats.sectorsWritten) * 0.5

		totalTPS += tps
		totalKBPerSec += kbPerSec

		diskStatsList = append(diskStatsList, &pb.DiskStats{
			Device:    stats.device,
			Tps:       tps,
			KbPerSec:  kbPerSec,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read /proc/diskstats: %w", err)
	}

	return &pb.DiskMetrics{
		Tps:      totalTPS,
		KbPerSec: totalKBPerSec,
		Disks:    diskStatsList,
	}, nil
}

