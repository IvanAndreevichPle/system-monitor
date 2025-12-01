package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// CollectNetworkStats reads /proc/net/dev and collects network statistics
// Returns network statistics with protocol and traffic information
// Note: For detailed protocol stats and top talkers, we need additional sources
// This implementation provides basic network interface statistics
func (c *Collector) CollectNetworkStats() (*pb.NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/net/dev: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var totalBytes uint64
	interfaceBytes := make(map[string]uint64)

	scanner := bufio.NewScanner(file)
	// Skip header lines
	for i := 0; i < 2 && scanner.Scan(); i++ {
		// Skip header
	}

	for scanner.Scan() {
		line := scanner.Text()
		// Format: "interface: bytes_rx packets_rx ... bytes_tx packets_tx ..."
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		interfaceName := strings.TrimSpace(parts[0])
		stats := strings.Fields(parts[1])
		if len(stats) < 16 {
			continue
		}

		// bytes_rx is at index 0, bytes_tx is at index 8
		bytesRx, err := strconv.ParseUint(stats[0], 10, 64)
		if err != nil {
			continue
		}

		bytesTx, err := strconv.ParseUint(stats[8], 10, 64)
		if err != nil {
			continue
		}

		totalBytesForInterface := bytesRx + bytesTx
		interfaceBytes[interfaceName] = totalBytesForInterface
		totalBytes += totalBytesForInterface
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read /proc/net/dev: %w", err)
	}

	// For now, return empty protocol and traffic stats
	// These will be populated in a future enhancement using /proc/net/sockstat
	// and netstat/ss for top talkers
	var protocolStats []*pb.ProtocolStats
	var trafficStats []*pb.TrafficStats

	// Calculate protocol percentages based on interface types
	// This is a simplified approach - in a full implementation,
	// we would parse /proc/net/sockstat for protocol-specific stats
	for iface, bytes := range interfaceBytes {
		if totalBytes > 0 {
			percent := float64(bytes) * 100.0 / float64(totalBytes)
			// Determine protocol based on interface name (simplified)
			protocol := "UNKNOWN"
			if strings.HasPrefix(iface, "eth") || strings.HasPrefix(iface, "en") {
				protocol = "ETHERNET"
			} else if iface == "lo" {
				protocol = "LOOPBACK"
			}

			protocolStats = append(protocolStats, &pb.ProtocolStats{
				Protocol: protocol,
				Bytes:    int64(bytes),
				Percent:  percent,
			})
		}
	}

	return &pb.NetworkStats{
		ProtocolStats: protocolStats,
		TrafficStats:  trafficStats,
	}, nil
}

