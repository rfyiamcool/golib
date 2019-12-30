package ringstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Normal(t *testing.T) {
	size := 10
	q := NewRingStore(size)
	for i := 0; i < 100; i++ {
		q.Add(i)
	}

	fmt.Println("len: ", q.Len())
	fmt.Println("source buf: ", q.buf)
	fmt.Println("get items: ", q.GetItems())
	fmt.Println("get item count: ", q.GetLimitItems(3))

	assert.Equal(t, size, q.Len())
	assert.Equal(t, q.GetLimitItems(3), []interface{}{97, 98, 99})
}

func Test_NotFull(t *testing.T) {
	size := 10
	q := NewRingStore(size)
	for i := 0; i < 3; i++ {
		q.Add(i)
	}

	fmt.Println("len: ", q.Len())
	fmt.Println("source buf: ", q.buf)
	fmt.Println("get items: ", q.GetItems())
	fmt.Println("get item count: ", q.GetLimitItems(5))

	assert.Equal(t, 3, q.Len())

	// dup test
	assert.Equal(t, q.GetLimitItems(5), []interface{}{0, 1, 2})
	assert.Equal(t, q.GetLimitItems(5), []interface{}{0, 1, 2})
}

func Test_Empty(t *testing.T) {
	size := 10
	q := NewRingStore(size)

	fmt.Println("len: ", q.Len())
	fmt.Println("get items: ", q.GetItems())
	fmt.Println("get item count: ", q.GetLimitItems(5))

	assert.Equal(t, 0, q.Len())
	assert.Equal(t, 0, len(q.GetLimitItems(5)))
}
