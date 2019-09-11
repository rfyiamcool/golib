package ringstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingStore(t *testing.T) {
	size := 10
	q := NewRingStore(size)
	for i := 0; i < 200000; i++ {
		q.Add(i)
	}

	fmt.Println("len: ", q.Len())
	fmt.Println("source buf: ", q.buf)
	fmt.Println("get items: ", q.GetItems())

	assert.Equal(t, size, q.Len())

	for i := 0; i < 200000; i++ {
		q.Add(i)
	}

	assert.Equal(t, size, q.Len())
}
