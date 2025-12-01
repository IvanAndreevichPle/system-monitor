package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectAll(t *testing.T) {
	// Skip if /proc doesn't exist (e.g., on non-Linux systems)
	if _, err := os.Stat("/proc/loadavg"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc filesystem not available")
	}

	c := NewCollector()
	snapshot, err := c.CollectAll()

	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.NotNil(t, snapshot.LoadAverage)
	assert.NotNil(t, snapshot.Cpu)
	assert.NotNil(t, snapshot.Disk)
	assert.NotNil(t, snapshot.Filesystem)
	assert.NotNil(t, snapshot.Network)
	assert.NotNil(t, snapshot.Sockets)
}

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	assert.NotNil(t, c)
}

