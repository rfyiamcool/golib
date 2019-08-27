package sema

import (
	"time"
)

type Semaphore struct {
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
		return true
	}

	tm := time.NewTimer(sem.timeout)
	defer tm.Stop()

	select {
	case <-sem.slots:
		return true
	case <-tm.C:
		return false
	}
}

func (sem *Semaphore) TryAcquire() bool {
	select {
	case <-sem.slots:
		return true
	default:
		return false
	}
}

func (sem *Semaphore) Release() {
	select {
	case sem.slots <- struct{}{}:
	default:
		panic("can not release")
	}
}

func (sem *Semaphore) SpareSem() int {
	return len(sem.slots)
}
