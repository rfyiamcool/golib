package set

import (
	"sync"
)

type Intset struct {
	needLock bool
	sync.RWMutex

	data map[int]bool
}

func NewIntset() *Intset {
	return &Intset{
		data: make(map[int]bool),
	}
}

func NewIntsetSafe() *Intset {
	return &Intset{
		data:     make(map[int]bool),
		needLock: true,
	}
}

func (s *Intset) Add(value int) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	s.data[value] = true
}

func (s *Intset) Remove(value int) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	delete(s.data, value)
}

func (s *Intset) Exists(value int) (existed bool) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	_, existed = s.data[value]
	return
}

func (s *Intset) Length() int {
	// don't need lock
	return len(s.data)
}

func (s *Intset) Values() (values []int) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	for val := range s.data {
		values = append(values, val)
	}

	return
}
