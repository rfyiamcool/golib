package anyhash

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	require.Equal(t, uint64(1), Hash(uint64(1)))
	require.Equal(t, uint64(1), Hash(1))
	require.Equal(t, uint64(2), Hash(int32(2)))
	require.Equal(t, uint64(math.MaxUint64)-1, Hash(int32(-2)))
	require.Equal(t, uint64(math.MaxUint64)-1, Hash(int64(-2)))
	require.Equal(t, uint64(3), Hash(uint32(3)))
	require.Equal(t, uint64(3), Hash(int64(3)))
}
