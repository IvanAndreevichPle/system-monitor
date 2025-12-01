package collector

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// CollectFilesystemInfo reads /proc/mounts and collects filesystem statistics
// Returns filesystem information with usage and inode statistics
func (c *Collector) CollectFilesystemInfo() (*pb.FilesystemInfo, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/mounts: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var filesystems []*pb.FilesystemStats

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		mountPoint := fields[1]

		// Skip special filesystems
		if strings.HasPrefix(mountPoint, "/proc") ||
			strings.HasPrefix(mountPoint, "/sys") ||
			strings.HasPrefix(mountPoint, "/dev") ||
			strings.HasPrefix(mountPoint, "/run") ||
			mountPoint == "/" && len(fields) > 2 && fields[2] == "rootfs" {
			continue
		}

		// Get filesystem statistics using statfs
		var stat unix.Statfs_t
		if err := unix.Statfs(mountPoint, &stat); err != nil {
			// Skip filesystems we can't stat
			continue
		}

		// Calculate used space
		totalBlocks := stat.Blocks
		freeBlocks := stat.Bavail
		usedBlocks := totalBlocks - freeBlocks

		// Block size is typically 4096 bytes (4 KB)
		blockSize := uint64(stat.Bsize)
		usedBytes := usedBlocks * blockSize
		usedMB := int64(usedBytes / (1024 * 1024))

		// Calculate used percentage
		var usedPercent float64
		if totalBlocks > 0 {
			usedPercent = float64(usedBlocks) * 100.0 / float64(totalBlocks)
		}

		// Calculate inode usage
		totalInodes := stat.Files
		freeInodes := stat.Ffree
		usedInodes := int64(totalInodes - freeInodes)

		// Calculate inode percentage
		var inodesPercent float64
		if totalInodes > 0 {
			inodesPercent = float64(usedInodes) * 100.0 / float64(totalInodes)
		}

		filesystems = append(filesystems, &pb.FilesystemStats{
			MountPoint:    mountPoint,
			UsedMb:        usedMB,
			UsedPercent:   usedPercent,
			UsedInodes:    usedInodes,
			InodesPercent: inodesPercent,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read /proc/mounts: %w", err)
	}

	return &pb.FilesystemInfo{
		Filesystems: filesystems,
	}, nil
}

