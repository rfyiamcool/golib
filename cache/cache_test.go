package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testResult struct {
	Code    int32
	Message string
	Result  interface{}
}

func TestCacheSimple(t *testing.T) {
	resp := &testResult{
		Code:    0,
		Message: "success",
		Result:  map[string]interface{}{"key1": "val1"},
	}
	key := "func:GetStats"
	SetWithExpire(key, resp, "1s")

	val, expired := Get(key)
	assert.Equal(t, expired, false)
	assert.NotEqual(t, val, nil)
	assert.Equal(t, val.(*testResult), resp)

	time.Sleep(3 * time.Second)
	val, expired = Get(key)
	t.Log("bad", val, expired)

	assert.Equal(t, expired, true)
	assert.NotEqual(t, val, nil)
	assert.Equal(t, val.(*testResult), resp)
	_ = val.(*testResult)

	if val != nil && !expired {
		t.Log(val.(*testResult))
	} else {
		t.Log("not data")
	}

	val = GetMust(key)
	assert.Equal(t, val, nil)
}

func TestCacheFetch(t *testing.T) {
	var (
		kname = "batch:get"
		ts    = time.Duration(10 * time.Second)
		start = time.Now()
		cnt   = 3
	)

	wg := sync.WaitGroup{}
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		idx := i
		go func() {
			defer wg.Done()
			v, err := GetSetWithLock(kname, time.Duration(5*time.Second), ts, func() (interface{}, error) {
				time.Sleep(3 * time.Second)
				return idx, nil
			})
			t.Log(">>>", v, err)
		}()
	}
	wg.Wait()

	t.Log(Get(kname))
	t.Log(Get(kname))
	t.Log(Get(kname))

	assert.Less(t, time.Since(start).Seconds(), float64(4))
}

func TestCacheFetchTimeout(t *testing.T) {
	var (
		timeout = time.Duration(1 * time.Second)
		ttl     = timeout
	)
	_, err := GetSetWithLock("nima", timeout, ttl, func() (interface{}, error) {
		time.Sleep(3 * time.Second)
		return 111, nil
	})

	v, expired := Get("nima")
	assert.Equal(t, v, nil)
	assert.Equal(t, expired, false)
	assert.Equal(t, err, ErrFetchTimeout)
}
