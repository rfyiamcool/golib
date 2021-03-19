package singleton

import "sync"

var (
	mutex sync.RWMutex

	onces = make(map[interface{}]*sync.Once)
)

func Once(key interface{}, fn func()) {
	mutex.Lock()
	once, ok := onces[key]
	if !ok {
		once = new(sync.Once)
	}
	onces[key] = once
	mutex.Unlock()

	once.Do(fn)
}

func Cancel(key interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(onces, key)
}

func Replace(key interface{}, once *sync.Once) {
	mutex.Lock()
	defer mutex.Unlock()

	onces[key] = once
}
