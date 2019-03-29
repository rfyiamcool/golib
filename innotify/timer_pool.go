package innotify

import (
	"sync"
	"time"
)

var globalTimerPool timerPool

type timerPool struct {
	p sync.Pool
}

func (tp *timerPool) Get(d time.Duration) *time.Timer {
	if t, _ := tp.p.Get().(*time.Timer); t != nil {
		t.Reset(d)
		return t
	}

	return time.NewTimer(d)
}

func (tp *timerPool) Put(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}

	tp.p.Put(t)
}
