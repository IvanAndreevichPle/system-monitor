package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectLoadAverage(t *testing.T) {
	// Skip if /proc/loadavg doesn't exist (e.g., on non-Linux systems)
	if _, err := os.Stat("/proc/loadavg"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/loadavg not available")
	}

	c := NewCollector()
	loadAvg, err := c.CollectLoadAverage()

	require.NoError(t, err)
	assert.NotNil(t, loadAvg)
	assert.GreaterOrEqual(t, loadAvg.Load1, 0.0)
	assert.GreaterOrEqual(t, loadAvg.Load5, 0.0)
	assert.GreaterOrEqual(t, loadAvg.Load15, 0.0)
}

func TestCollectLoadAverage_InvalidFormat(t *testing.T) {
	// This test would require mocking /proc/loadavg
	// For now, we test with real file which should always be valid
	t.Skip("Skipping: requires mocking /proc/loadavg")
}

