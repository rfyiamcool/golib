package set

import (
	"sync"
)

type strset struct {
	needLock bool
	sync.RWMutex

	data map[string]bool
}

func NewStringSet() *strset {
	return &strset{
		data: make(map[string]bool),
	}
}

func NewStringSetSafe() *strset {
	return &strset{
		data:     make(map[string]bool),
		needLock: true,
	}
}

func (s *strset) Add(value string) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	s.data[value] = true
}

func (s *strset) Remove(value string) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	delete(s.data, value)
}

func (s *strset) Contains(value string) (exists bool) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	_, exists = s.data[value]
	return
}

func (s *strset) Length() int {
	// don't need lock
	return len(s.data)
}

func (s *strset) Values() (values []string) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	for val := range s.data {
		values = append(values, val)
	}

	return
}
