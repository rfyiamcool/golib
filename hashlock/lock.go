package hashlock

import (
	"hash/fnv"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type HashMutex struct {
	mutexes []sync.Mutex
}

func New(n int) *HashMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}

	return &HashMutex{
		mutexes: make([]sync.Mutex, n),
	}
}

func (mu *HashMutex) Lock(id string) {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].Lock()
}

func (mu *HashMutex) Unlock(id string) error {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].Unlock()
	return nil
}

type HashRWMutex struct {
	mutexes []sync.RWMutex
}

func NewMultiRWMutex(n int) *HashRWMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}

	return &HashRWMutex{
		mutexes: make([]sync.RWMutex, n),
	}
}

func (mu *HashRWMutex) Lock(id string) {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].Lock()
}

func (mu *HashRWMutex) Unlock(id string) error {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].Unlock()
	return nil
}

func (mu *HashRWMutex) RLock(id string) {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].RLock()
}

func (mu *HashRWMutex) RUnlock(id string) error {
	idx := mod(id, len(mu.mutexes))
	mu.mutexes[idx].RUnlock()
	return nil
}

var (
	mlocks = NewMultiRWMutex(512)
)

func Lock(id string) {
	mlocks.Lock(id)
}

func Unlock(id string) error {
	return mlocks.Unlock(id)
}

func RLock(id string) {
	mlocks.RLock(id)
}

func RUnLock(id string) error {
	return mlocks.RUnlock(id)
}

func mod(id string, length int) uint32 {
	return hash(id) % uint32(length)
}

func hash(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}

const mutexLocked = 1 << iota

func TryLock(m *sync.Mutex) bool {
	// 0, unlock
	// 1, lock
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m)), 0, mutexLocked)
}

func CheckLocked(m *sync.Mutex) bool {
	return atomic.LoadInt32((*int32)(unsafe.Pointer(m))) == int32(1)
}
