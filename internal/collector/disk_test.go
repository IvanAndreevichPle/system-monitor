package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectDiskMetrics(t *testing.T) {
	// Skip if /proc/diskstats doesn't exist (e.g., on non-Linux systems)
	if _, err := os.Stat("/proc/diskstats"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/diskstats not available")
	}

	c := NewCollector()
	disk, err := c.CollectDiskMetrics()

	require.NoError(t, err)
	assert.NotNil(t, disk)
	assert.GreaterOrEqual(t, disk.Tps, 0.0)
	assert.GreaterOrEqual(t, disk.KbPerSec, 0.0)
	// Should have at least some disks (or empty list)
	assert.NotNil(t, disk.Disks)
}

