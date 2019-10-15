package trylock

import (
	"sync"
	"testing"
)

func TestTryLock(t *testing.T) {
	var mu Mutex
	if !mu.TryLock() {
		t.Fatal("mutex must be unlocked")
	}
	if mu.TryLock() {
		t.Fatal("mutex must be locked")
	}

	mu.Unlock()
	if !mu.TryLock() {
		t.Fatal("mutex must be unlocked")
	}
	if mu.TryLock() {
		t.Fatal("mutex must be locked")
	}

	mu.Unlock()
	mu.Lock()
	if mu.TryLock() {
		t.Fatal("mutex must be locked")
	}
	if mu.TryLock() {
		t.Fatal("mutex must be locked")
	}
	mu.Unlock()
}

func TestBlock(t *testing.T) {
	var mu Mutex
	var x int
	var wg sync.WaitGroup

	for i := 0; i < 1024; i++ {
		wg.Add(1)
		if i%2 == 0 {
			go func() {
				defer wg.Done()
				if mu.TryLock() {
					x++
					mu.Unlock()
				}
			}()
			continue
		}
		go func() {
			defer wg.Done()
			mu.Lock()
			x++
			mu.Unlock()
		}()
	}
	wg.Wait()
	t.Log("finish")
}
