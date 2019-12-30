package ringstore

import (
	"sync"
)

// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

type RingStore struct {
	sync.RWMutex

	buf    []interface{}
	head   int
	tail   int
	count  int
	minCap int

	size int
}

func NewRingStore(size int) *RingStore {
	return &RingStore{
		size: size,
	}
}

func (q *RingStore) Len() int {
	return q.count
}

func (q *RingStore) GetLimitItems(n int) []interface{} {
	var (
		list []interface{}
		tail = q.prev(q.tail)
	)

	if len(q.buf) == 0 {
		return list
	}

	if n > q.count {
		n = q.count
	}

	for index := 0; index < n; index++ {
		ret := q.buf[tail]
		tail = q.prev(tail)
		list = append(list, ret)
	}

	return q.aesOrder(list)
}

func (q *RingStore) GetItems() []interface{} {
	var (
		list []interface{}
		head = q.head
	)

	if len(q.buf) == 0 {
		return list
	}

	for index := 0; index < q.size; index++ {
		ret := q.buf[head]
		head = q.next(head)
		list = append(list, ret)
	}

	return list
}

func (q *RingStore) Add(elem interface{}) {
	q.push(elem)
	if q.Len() > q.size {
		q.pop()
	}
}

// push back
func (q *RingStore) push(elem interface{}) {
	q.growIfFull()

	q.buf[q.tail] = elem
	q.tail = q.next(q.tail)
	q.count++
}

// pop front
func (q *RingStore) pop() interface{} {
	if q.count <= 0 {
		panic("RingStore: PopFront() called on empty queue")
	}

	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return ret
}

func (q *RingStore) Front() interface{} {
	if q.count <= 0 {
		panic("RingStore: Front() called when empty")
	}
	return q.buf[q.head]
}

func (q *RingStore) Back() interface{} {
	if q.count <= 0 {
		panic("RingStore: Back() called when empty")
	}
	return q.buf[q.prev(q.tail)]
}

func (q *RingStore) Clear() {
	modBits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & modBits {
		q.buf[h] = nil
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// prev returns the previous buffer position wrapping around buffer.
func (q *RingStore) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *RingStore) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

func (q *RingStore) growIfFull() {
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]interface{}, q.minCap)
		return
	}
	if q.count == len(q.buf) {
		q.resize()
	}
}

func (q *RingStore) shrinkIfExcess() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

func (q *RingStore) resize() {
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

func (q *RingStore) aesOrder(s []interface{}) []interface{} {
	for from, to := 0, len(s)-1; from < to; from, to = from+1, to-1 {
		s[from], s[to] = s[to], s[from]
	}

	return s
}
