package collector

import (
	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// Collector collects system metrics
type Collector struct{}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{}
}

// CollectAll collects all system metrics and returns a MetricsSnapshot
func (c *Collector) CollectAll() (*pb.MetricsSnapshot, error) {
	snapshot := &pb.MetricsSnapshot{}

	// Collect Load Average
	loadAvg, err := c.CollectLoadAverage()
	if err != nil {
		return nil, err
	}
	snapshot.LoadAverage = loadAvg

	// Collect CPU metrics
	cpuMetrics, err := c.CollectCPUMetrics()
	if err != nil {
		return nil, err
	}
	snapshot.Cpu = cpuMetrics

	// Collect disk metrics
	diskMetrics, err := c.CollectDiskMetrics()
	if err != nil {
		return nil, err
	}
	snapshot.Disk = diskMetrics

	// Collect filesystem info
	filesystemInfo, err := c.CollectFilesystemInfo()
	if err != nil {
		return nil, err
	}
	snapshot.Filesystem = filesystemInfo

	// Collect network stats
	networkStats, err := c.CollectNetworkStats()
	if err != nil {
		return nil, err
	}
	snapshot.Network = networkStats

	// Collect socket stats
	socketStats, err := c.CollectSocketStats()
	if err != nil {
		return nil, err
	}
	snapshot.Sockets = socketStats

	return snapshot, nil
}

