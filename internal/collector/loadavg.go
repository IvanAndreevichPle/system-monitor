package collector

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// CollectLoadAverage reads and parses /proc/loadavg
// Format: "0.52 0.58 0.59 1/396 12345"
// Returns LoadAverage with load1, load5, load15
func (c *Collector) CollectLoadAverage() (*pb.LoadAverage, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/loadavg: %w", err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid /proc/loadavg format: expected at least 3 fields, got %d", len(fields))
	}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse load1: %w", err)
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse load5: %w", err)
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse load15: %w", err)
	}

	return &pb.LoadAverage{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

