package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectCPUMetrics(t *testing.T) {
	// Skip if /proc/stat doesn't exist (e.g., on non-Linux systems)
	if _, err := os.Stat("/proc/stat"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/stat not available")
	}

	c := NewCollector()
	cpu, err := c.CollectCPUMetrics()

	require.NoError(t, err)
	assert.NotNil(t, cpu)
	assert.GreaterOrEqual(t, cpu.UserPercent, 0.0)
	assert.LessOrEqual(t, cpu.UserPercent, 100.0)
	assert.GreaterOrEqual(t, cpu.SystemPercent, 0.0)
	assert.LessOrEqual(t, cpu.SystemPercent, 100.0)
	assert.GreaterOrEqual(t, cpu.IdlePercent, 0.0)
	assert.LessOrEqual(t, cpu.IdlePercent, 100.0)

	// Sum should be approximately 100%
	total := cpu.UserPercent + cpu.SystemPercent + cpu.IdlePercent
	assert.InDelta(t, 100.0, total, 1.0) // Allow 1% tolerance
}

