package latency

import (
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
	block          bool
	lastInsertTime time.Time
}

func NewLatency(timeout time.Duration, block bool) *Latency {
	return &Latency{
		timeout: timeout,
		block:   block,
	}
}

func (l *Latency) Push(value time.Duration) {
	t := int64(value / time.Microsecond)
	if t <= 0 {
		return
	}

	if !l.Lock() {
		return
	}

	if time.Now().Sub(l.lastInsertTime) >= l.timeout {
		l.count = 0
		l.sum = 0
	}

	l.count++
	l.sum += t

	l.lastInsertTime = time.Now()
	l.Release()
}

func (l *Latency) Reset() {
	l.count = 0
	l.sum = 0
}

func (l *Latency) Lock() bool {
	if !l.block {
		return atomic.CompareAndSwapInt32(&l.flag, free, working)
	}

	for {
		if atomic.CompareAndSwapInt32(&l.flag, free, working) {
			return true
		}
	}
}

func (l *Latency) Release() {
	atomic.StoreInt32(&l.flag, free)
}

func (l *Latency) Calc() time.Duration {
	if time.Now().Sub(l.lastInsertTime) >= l.timeout {
		return 0
	}

	avg := l.sum / l.count
	return time.Duration(avg) * time.Microsecond
}
