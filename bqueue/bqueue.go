package bqueue

import (
	"container/list"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Queue interface {
	Put(item interface{})

	PutTimeout(item interface{}, timeout time.Duration) bool

	Poll() interface{}

	PollTimeout(timeout time.Duration) (interface{}, bool)

	Len() int
}

func NewQueue(capacity int) Queue {
	if capacity <= 0 {
		c := make(chan struct{})
		return &infiniteQueue{
			store: list.New(),
			empty: unsafe.Pointer(&c),
		}
	}

	return &finiteQueue{
		store: make(chan interface{}, capacity),
	}
}

type infiniteQueue struct {
	sync.Mutex
	store *list.List
	empty unsafe.Pointer
}

var _ Queue = &infiniteQueue{}

func (q *infiniteQueue) Put(item interface{}) {
	if isNil(item) {
		return
	}

	q.Lock()
	defer q.Unlock()

	q.store.PushBack(item)
	if q.store.Len() < 2 {
		// empty -> has one element
		q.broadcast()
	}
}

func (q *infiniteQueue) PutTimeout(item interface{}, timeout time.Duration) bool {
	q.Put(item)
	return !isNil(item)
}

func (q *infiniteQueue) Poll() interface{} {
	q.Lock()
	defer q.Unlock()

	for q.store.Len() == 0 {
		q.wait()
	}
	item := q.store.Front()
	q.store.Remove(item)
	return item.Value
}

func (q *infiniteQueue) PollTimeout(timeout time.Duration) (interface{}, bool) {
	deadline := time.Now().Add(timeout)
	q.Lock()
	defer q.Unlock()

	for q.store.Len() == 0 {
		timeout = -time.Since(deadline)
		if timeout <= 0 || !q.waitTimeout(timeout) {
			return nil, false
		}
	}

	item := q.store.Front()
	q.store.Remove(item)
	return item.Value, true
}

func (q *infiniteQueue) Len() int {
	q.Lock()
	defer q.Unlock()
	return q.store.Len()
}

func (q *infiniteQueue) wait() {
	c := q.notifyChan()
	q.Unlock()
	defer q.Lock()
	<-c
}

func (q *infiniteQueue) waitTimeout(timeout time.Duration) bool {
	c := q.notifyChan()
	q.Unlock()
	defer q.Lock()

	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (q *infiniteQueue) notifyChan() <-chan struct{} {
	ptr := atomic.LoadPointer(&q.empty)
	return *((*chan struct{})(ptr))
}

func (q *infiniteQueue) broadcast() {
	c := make(chan struct{})
	old := atomic.SwapPointer(&q.empty, unsafe.Pointer(&c))
	close(*(*chan struct{})(old))
}

type finiteQueue struct {
	store chan interface{}
}

func (q *finiteQueue) Put(item interface{}) {
	if isNil(item) {
		return
	}
	q.store <- item
}

func (q *finiteQueue) PutTimeout(item interface{}, timeout time.Duration) bool {
	if isNil(item) {
		return false
	}

	if timeout <= 0 {
		select {
		case q.store <- item:
			return true
		default:
			return false
		}
	}

	select {
	case q.store <- item:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (q *finiteQueue) Poll() interface{} {
	item := <-q.store
	return item
}

func (q *finiteQueue) PollTimeout(timeout time.Duration) (interface{}, bool) {
	if timeout <= 0 {
		select {
		case item := <-q.store:
			return item, true
		default:
			return nil, false
		}
	}

	select {
	case item := <-q.store:
		return item, true
	case <-time.After(timeout):
		return nil, false
	}
}

func (q *finiteQueue) Len() int {
	return len(q.store)
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}
