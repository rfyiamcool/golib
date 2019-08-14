package bqueue

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	q := NewQueue(0)

	q.Put(nil)
	assert.Equal(t, q.Len(), 0)

	q.PutTimeout(nil, 0)
	assert.Equal(t, q.Len(), 0)

	q.Put(0)
	assert.Equal(t, q.Len(), 1)

	v := q.Poll()
	assert.Equal(t, v.(int), 0)
}

func Test_PollTimeout(t *testing.T) {
	q := NewQueue(10)
	wg := sync.WaitGroup{}

	var (
		cnt     int32
		idx     int32
		loopNum int32 = 50
	)
	go func() {
		wg.Add(1)
		for idx < int32(loopNum) {
			idx++
			q.Poll()
			atomic.AddInt32(&cnt, 1)
		}
		wg.Done()
	}()

	var index int32
	for index = 0; index < loopNum; index++ {
		q.Put(index)
	}

	wg.Wait()
	t.Log(">>> ", cnt, idx)
	assert.Equal(t, cnt, loopNum)
	assert.Equal(t, idx, loopNum)
}
