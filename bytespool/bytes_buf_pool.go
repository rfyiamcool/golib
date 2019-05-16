package bytesbuf

import "sync"


const (
	base      = 64
	sizeMulti = 2

	numPools  = 5
)

var (
	pool     [numPools]sync.Pool
	poolSize [numPools]int32
)

func createAllocFunc(size int32) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

func init() {
	size := int32(base)
	for i := 0; i < numPools; i++ {
		pool[i] = sync.Pool{
			New: createAllocFunc(size),
		}
		poolSize[i] = size
		size *= sizeMulti
	}
}

func GetPool(size int32) *sync.Pool {
	for idx, ps := range poolSize {
		if size <= ps {
			return &pool[idx]
		}
	}

	return nil
}

func Alloc(size int32) []byte {
	pool := GetPool(size)
	if pool != nil {
		return pool.Get().([]byte)
	}

	return make([]byte, size)
}

func Free(b []byte) {
	size := int32(cap(b))
	b = b[0:cap(b)]
	for i := numPools - 1; i >= 0; i-- {
		if size >= poolSize[i] {
			pool[i].Put(b)
			return
		}
	}
}
