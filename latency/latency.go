package latency

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	free = iota
	working
)

type Latency struct {
	timeout        time.Duration
	count          int64
	sum            int64
	flag           int32
	lastInsertTime time.Time
}

func NewLatency(timeout time.Duration) *Latency {
	return &Latency{
		timeout: timeout,
	}
}

func (l *Latency) Push(value time.Duration) {
	t := int64(value / time.Microsecond)
	if t <= 0 {
		return
	}

	if !atomic.CompareAndSwapInt32(&l.flag, free, working) {
		return
	}

	if time.Now().Sub(l.lastInsertTime) >= l.timeout {
		l.reset()
	}

	l.count++
	l.sum += t

	l.lastInsertTime = time.Now()
	atomic.StoreInt32(&l.flag, free)
}

func (l *Latency) reset() {
	l.count = 0
	l.sum = 0
}

func (l *Latency) Calc() time.Duration {
	if time.Now().Sub(l.lastInsertTime) >= l.timeout {
		return 0
	}

	avg := l.sum / l.count
	return time.Duration(avg)
}
