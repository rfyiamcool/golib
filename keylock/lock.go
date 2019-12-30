package keylock

import (
	"hash/fnv"
	"runtime"
	"sync"
)

type MultiMutex struct {
	mutexes []sync.Mutex
}

func New(n int) *MultiMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}

	return &MultiMutex{
		mutexes: make([]sync.Mutex, n),
	}
}

func (mu *MultiMutex) Lock(id string) {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].Lock()
}

func (mu *MultiMutex) Unlock(id string) error {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].Unlock()
	return nil
}

type MultiRWMutex struct {
	mutexes []sync.RWMutex
}

func NewMultiRWMutex(n int) *MultiRWMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}

	return &MultiRWMutex{
		mutexes: make([]sync.RWMutex, n),
	}
}

func (mu *MultiRWMutex) Lock(id string) {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].Lock()
}

func (mu *MultiRWMutex) Unlock(id string) error {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].Unlock()
	return nil
}

func (mu *MultiRWMutex) RLock(id string) {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].RLock()
}

func (mu *MultiRWMutex) RUnlock(id string) error {
	mu.mutexes[hash(id)%uint32(len(mu.mutexes))].RUnlock()
	return nil
}

func hash(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
