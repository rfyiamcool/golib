package sema

import (
	"time"
)

type Semaphore struct {
	locked  bool
	slots   chan struct{}
	timeout time.Duration
}

func NewSemaphore(count int, timeout time.Duration) *Semaphore {
	sem := &Semaphore{
		slots:   make(chan struct{}, count),
		timeout: timeout,
	}

	for i := 0; i < count; i++ {
		sem.slots <- struct{}{}
	}

	return sem
}

func (sem *Semaphore) Acquire() bool {
	if sem.timeout == 0 {
		<-sem.slots
		sem.locked = true
		return true
	}

	tm := time.NewTimer(sem.timeout)
	defer tm.Stop()

	select {
	case <-sem.slots:
		sem.locked = true
		return true
	case <-tm.C:
		return false
	}
}

func (sem *Semaphore) TryAcquire() bool {
	select {
	case <-sem.slots:
		sem.locked = true
		return true
	default:
		return false
	}
}

func (sem *Semaphore) Release() {
	if sem.locked {
		sem.slots <- struct{}{}
		sem.locked = false
		return
	}

	panic("can not release")
}

func (sem *Semaphore) Size() int {
	return len(sem.slots)
}

func (sem *Semaphore) IsLocked() bool {
	return sem.locked
}
