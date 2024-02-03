package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedDiscoverer(t *testing.T) {
	d := FixedDiscoverer([]string{"a", "b"})
	instances, err := d.Instances()
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, instances)
}
