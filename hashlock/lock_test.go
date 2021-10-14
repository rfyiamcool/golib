package hashlock

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryLock(t *testing.T) {
	m := new(sync.Mutex)
	b := CheckLocked(m)

	m.Lock()
	b = CheckLocked(m)
	assert.Equal(t, true, b)

	m.Unlock()
	b = CheckLocked(m)
	assert.Equal(t, false, b)

	for i := 0; i < 10; i++ {
		b = CheckLocked(m)
		assert.Equal(t, false, b)
	}

	b = TryLock(m)
	assert.Equal(t, true, b)

	m.Unlock()
	b = CheckLocked(m)
	assert.Equal(t, false, b)
}
