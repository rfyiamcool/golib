package atomic2

import (
	"sync"
	"sync/atomic"
	"time"
)

type AtomicInt32 struct {
	int32
}

func NewAtomicInt32(n int32) AtomicInt32 {
	return AtomicInt32{n}
}

func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32(&i.int32, n)
}

func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32(&i.int32, n)
}

func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&i.int32)
}

func (i *AtomicInt32) CompareAndSwap(oldval, newval int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&i.int32, oldval, newval)
}

type AtomicInt64 struct {
	int64
}

func NewAtomicInt64(n int64) AtomicInt64 {
	return AtomicInt64{n}
}

func (i *AtomicInt64) Add(n int64) int64 {
	return atomic.AddInt64(&i.int64, n)
}

func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64(&i.int64, n)
}

func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&i.int64)
}

func (i *AtomicInt64) CompareAndSwap(oldval, newval int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&i.int64, oldval, newval)
}

type AtomicDuration struct {
	int64
}

func NewAtomicDuration(duration time.Duration) AtomicDuration {
	return AtomicDuration{int64(duration)}
}

func (d *AtomicDuration) Add(duration time.Duration) time.Duration {
	return time.Duration(atomic.AddInt64(&d.int64, int64(duration)))
}

func (d *AtomicDuration) Set(duration time.Duration) {
	atomic.StoreInt64(&d.int64, int64(duration))
}

func (d *AtomicDuration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64(&d.int64))
}

func (d *AtomicDuration) CompareAndSwap(oldval, newval time.Duration) (swapped bool) {
	return atomic.CompareAndSwapInt64(&d.int64, int64(oldval), int64(newval))
}

type AtomicBool struct {
	int32
}

func NewAtomicBool(n bool) AtomicBool {
	if n {
		return AtomicBool{1}
	}
	return AtomicBool{0}
}

func (i *AtomicBool) Set(n bool) {
	if n {
		atomic.StoreInt32(&i.int32, 1)
	} else {
		atomic.StoreInt32(&i.int32, 0)
	}
}

func (i *AtomicBool) Get() bool {
	return atomic.LoadInt32(&i.int32) != 0
}

func (i *AtomicBool) CompareAndSwap(o, n bool) bool {
	var old, new int32
	if o {
		old = 1
	}
	if n {
		new = 1
	}
	return atomic.CompareAndSwapInt32(&i.int32, old, new)
}

type AtomicString struct {
	mu  sync.Mutex
	str string
}

func (s *AtomicString) Set(str string) {
	s.mu.Lock()
	s.str = str
	s.mu.Unlock()
}

func (s *AtomicString) Get() string {
	s.mu.Lock()
	str := s.str
	s.mu.Unlock()
	return str
}

func (s *AtomicString) CompareAndSwap(oldval, newval string) (swqpped bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.str == oldval {
		s.str = newval
		return true
	}
	return false
}
