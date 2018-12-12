package timerPool

import (
	"sync"
	"time"
)

// global pool of *time.Timer's. can be used by multiple goroutines concurrently.
var GlobalTimerPool = NewTimerPool()

// NewTImerPool create timePool obj
func NewTimerPool() *timerPool {
	return &timerPool{}
}

// timerPool provides GC-able pooling of *time.Timer's.
// can be used by multiple goroutines concurrently.
type timerPool struct {
	p sync.Pool
}

// Get returns a timer that completes after the given duration.
func (tp *timerPool) Get(d time.Duration) *time.Timer {
	if t, _ := tp.p.Get().(*time.Timer); t != nil {
		t.Reset(d)
		return t
	}

	return time.NewTimer(d)
}

// Put pools the given timer.
func (tp *timerPool) Put(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}

	tp.p.Put(t)
}
