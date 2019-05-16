package set

import (
	"sync"
)

type set struct {
	needLock bool
	sync.RWMutex

	data map[string]bool
}

func NewSet() *set {
	return &set{
		data: make(map[string]bool),
	}
}

func NewSetSafe() *set {
	return &set{
		data:     make(map[string]bool),
		needLock: true,
	}
}

func (s *set) Add(value string) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	s.data[value] = true
}

func (s *set) Remove(value string) {
	if s.needLock {
		s.Lock()
		defer s.Unlock()
	}

	delete(s.data, value)
}

func (s *set) Contains(value string) (exists bool) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	_, exists = s.data[value]
	return
}

func (s *set) Length() int {
	// don't need lock
	return len(s.data)
}

func (s *set) Values() (values []string) {
	if s.needLock {
		s.RLock()
		defer s.RUnlock()
	}

	for val := range s.data {
		values = append(values, val)
	}

	return
}
