package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// cpuStats holds raw CPU statistics from /proc/stat
type cpuStats struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
	guest   uint64
	guestNice uint64
}

// parseCPUStats parses the first line of /proc/stat (aggregate CPU stats)
func parseCPUStats(line string) (*cpuStats, error) {
	fields := strings.Fields(line)
	if len(fields) < 11 {
		return nil, fmt.Errorf("invalid /proc/stat format: expected at least 11 fields, got %d", len(fields))
	}

	// Skip "cpu" field
	values := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		val, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CPU stat field %d: %w", i, err)
		}
		values[i] = val
	}

	return &cpuStats{
		user:      values[0],
		nice:      values[1],
		system:    values[2],
		idle:      values[3],
		iowait:    values[4],
		irq:       values[5],
		softirq:   values[6],
		steal:     values[7],
		guest:     values[8],
		guestNice: values[9],
	}, nil
}

// CollectCPUMetrics reads and parses /proc/stat
// Returns CPU metrics with percentages
// Note: For accurate percentages, we need two samples. This implementation
// returns current values. Averaging will be implemented in stage 10.
func (c *Collector) CollectCPUMetrics() (*pb.CPUMetrics, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/stat: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read /proc/stat: empty file")
	}

	stats, err := parseCPUStats(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Calculate total time
	// Note: guest time is already included in user/nice, so we subtract it
	userTime := stats.user - stats.guest
	niceTime := stats.nice - stats.guestNice
	totalTime := userTime + niceTime + stats.system + stats.idle + stats.iowait +
		stats.irq + stats.softirq + stats.steal

	if totalTime == 0 {
		return &pb.CPUMetrics{
			UserPercent:   0,
			SystemPercent: 0,
			IdlePercent:   0,
		}, nil
	}

	// Calculate percentages
	userPercent := float64(userTime+niceTime) * 100.0 / float64(totalTime)
	systemPercent := float64(stats.system+stats.irq+stats.softirq) * 100.0 / float64(totalTime)
	idlePercent := float64(stats.idle+stats.iowait) * 100.0 / float64(totalTime)

	return &pb.CPUMetrics{
		UserPercent:   userPercent,
		SystemPercent: systemPercent,
		IdlePercent:   idlePercent,
	}, nil
}

