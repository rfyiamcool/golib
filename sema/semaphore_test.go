package sema

import (
	"testing"
	"time"
)

func TestSemaNoTimeout(t *testing.T) {
	s := NewSemaphore(1, 0)
	s.Acquire()
	released := false

	go func() {
		time.Sleep(10 * time.Millisecond)
		released = true
		s.Release()
	}()

	s.Acquire()
	if !released {
		t.Errorf("release: false, want true")
	}
}

func TestSemaTimeout(t *testing.T) {
	s := NewSemaphore(1, 5*time.Millisecond)
	s.Acquire()

	go func() {
		time.Sleep(10 * time.Millisecond)
		s.Release()
	}()

	if s.Acquire() {
		t.Errorf("Acquire: true, want false")
	}
}
