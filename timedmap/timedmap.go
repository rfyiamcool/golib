package timedmap

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type callback func(value interface{})

type TimedMap struct {
	mtx             sync.RWMutex
	container       map[keyWrap]*element
	cleanupTickTime time.Duration
	cleaner         *time.Ticker
	cleanerStopChan chan bool
}

type keyWrap struct {
	sec int
	key interface{}
}

type element struct {
	value   interface{}
	expires time.Time
	cbs     []callback
}

func New(cleanupTickTime time.Duration) *TimedMap {
	tm := &TimedMap{
		container:       make(map[keyWrap]*element),
		cleanerStopChan: make(chan bool),
	}

	tm.cleaner = time.NewTicker(cleanupTickTime)

	go func() {
		for {
			select {
			case <-tm.cleaner.C:
				tm.cleanUp()
			case <-tm.cleanerStopChan:
				break
			}
		}
	}()

	return tm
}

func (tm *TimedMap) Section(i int) Section {
	return newSection(tm, i)
}

func (tm *TimedMap) Ident() int {
	return 0
}

func (tm *TimedMap) Set(key, value interface{}, expiresAfter time.Duration, cb ...callback) {
	tm.set(key, 0, value, expiresAfter, cb...)
}

func (tm *TimedMap) GetValue(key interface{}) interface{} {
	v := tm.get(key, 0)
	if v == nil {
		return nil
	}
	return v.value
}

func (tm *TimedMap) GetExpires(key interface{}) (time.Time, error) {
	v := tm.get(key, 0)
	if v == nil {
		return time.Time{}, ErrKeyNotFound
	}
	return v.expires, nil
}

func (tm *TimedMap) SetExpire(key interface{}, d time.Duration) error {
	return tm.setExpire(key, 0, d)
}

func (tm *TimedMap) Contains(key interface{}) bool {
	return tm.get(key, 0) != nil
}

func (tm *TimedMap) Remove(key interface{}) {
	tm.remove(key, 0)
}

func (tm *TimedMap) Refresh(key interface{}, d time.Duration) error {
	return tm.refresh(key, 0, d)
}

func (tm *TimedMap) Flush() {
	tm.container = make(map[keyWrap]*element)
}

func (tm *TimedMap) Size() int {
	return len(tm.container)
}

func (tm *TimedMap) StopCleaner() {
	go func() {
		tm.cleanerStopChan <- true
	}()
	tm.cleaner.Stop()
}

func (tm *TimedMap) expireElement(key interface{}, sec int, v *element) {
	for _, cb := range v.cbs {
		cb(v.value)
	}

	k := keyWrap{
		sec: sec,
		key: key,
	}

	delete(tm.container, k)
}

func (tm *TimedMap) cleanUp() {
	now := time.Now()

	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	for k, v := range tm.container {
		if now.After(v.expires) {
			tm.expireElement(k.key, k.sec, v)
		}
	}
}

func (tm *TimedMap) set(key interface{}, sec int, val interface{}, expiresAfter time.Duration, cb ...callback) {
	// re-use element when existent on this key
	if v := tm.getRaw(key, sec); v != nil {
		v.value = val
		v.expires = time.Now().Add(expiresAfter)
		v.cbs = cb
		return
	}

	k := keyWrap{
		sec: sec,
		key: key,
	}

	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	tm.container[k] = &element{
		value:   val,
		expires: time.Now().Add(expiresAfter),
		cbs:     cb,
	}
}

func (tm *TimedMap) get(key interface{}, sec int) *element {
	v := tm.getRaw(key, sec)

	if v == nil {
		return nil
	}

	if time.Now().After(v.expires) {
		tm.expireElement(key, sec, v)
		return nil
	}

	return v
}

func (tm *TimedMap) getRaw(key interface{}, sec int) *element {
	k := keyWrap{
		sec: sec,
		key: key,
	}

	tm.mtx.RLock()
	v, ok := tm.container[k]
	tm.mtx.RUnlock()

	if !ok {
		return nil
	}

	return v
}

func (tm *TimedMap) remove(key interface{}, sec int) {
	k := keyWrap{
		sec: sec,
		key: key,
	}

	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	delete(tm.container, k)
}

func (tm *TimedMap) refresh(key interface{}, sec int, d time.Duration) error {
	v := tm.get(key, sec)
	if v == nil {
		return ErrKeyNotFound
	}
	v.expires = v.expires.Add(d)
	return nil
}

func (tm *TimedMap) setExpire(key interface{}, sec int, d time.Duration) error {
	v := tm.get(key, sec)
	if v == nil {
		return ErrKeyNotFound
	}
	v.expires = time.Now().Add(d)
	return nil
}
