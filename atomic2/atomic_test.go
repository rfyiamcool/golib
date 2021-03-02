package atomic2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	state := NewAtomicBool(false)
	state.SetTrue()
	assert.Equal(t, true, state.IsTrue())

	state.SetFalse()
	assert.Equal(t, false, state.IsTrue())

	state.Set(true)
	assert.Equal(t, true, state.IsTrue())
	assert.Equal(t, true, state.Get())

	state.CompareAndSwap(true, false)
	assert.Equal(t, false, state.Get())
	assert.Equal(t, true, state.IsFalse()) // is false
}
