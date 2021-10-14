package ringqueue

import (
	"sync"
)

const defaultMinQueueLen = 16

type RingQueue struct {
	sync.RWMutex
	size int

	buf               []interface{}
	head, tail, count int
}

func New() *RingQueue {
	return &RingQueue{
		buf: make([]interface{}, defaultMinQueueLen),
	}
}

func NewWithOption(size int) *RingQueue {
	return &RingQueue{
		buf:  make([]interface{}, size),
		size: size,
	}
}

func (q *RingQueue) Length() int {
	return q.count
}

func (q *RingQueue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *RingQueue) Add(elem interface{}) {
	q.Lock()
	defer q.Unlock()

	if q.count == len(q.buf) {
		q.resize()
	}

	q.count++
	q.buf[q.tail] = elem

	if q.tail+1 < len(q.buf) {
		q.tail++
	}
	if len(q.buf) == q.count {
		q.tail = 0
	}
}

func (q *RingQueue) Peek() interface{} {
	q.RLock()
	defer q.RUnlock()

	if q.count <= 0 {
		return nil
	}
	return q.buf[q.head]
}

func (q *RingQueue) Get(i int) interface{} {
	q.RLock()
	defer q.RUnlock()

	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		return nil
	}

	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

func (q *RingQueue) Remove() interface{} {
	q.Lock()
	defer q.Unlock()

	if q.count <= 0 {
		return nil
	}

	ret := q.buf[q.head]
	q.buf[q.head] = nil

	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--

	// Resize down if buffer 1/4 full.
	if len(q.buf) > q.size && (q.count<<2) == len(q.buf) {
		q.resize()
	}

	return ret
}
