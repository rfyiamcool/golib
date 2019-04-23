package latency

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	n := NewLatency(
		10*time.Second,
		true,
	)

	n.Push(1 * time.Second)
	n.Push(1 * time.Second)

	assert.Equal(t, n.Calc(), 1*time.Second)
}

func TestNonBlock(t *testing.T) {
	n := NewLatency(
		10*time.Second,
		false,
	)

	n.Push(1 * time.Second)
	n.Push(1 * time.Second)

	assert.Equal(t, n.Calc(), 1*time.Second)
}
