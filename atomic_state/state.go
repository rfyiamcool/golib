package atomicState

import (
	"sync/atomic"
)

var (
	SetFlag   int32 = 1
	UnsetFlag int32 = 0
)

func SetFlag(n *int32) bool {
	return atomic.CompareAndSwapInt32(n, unsetFlag, setFlag)
}

func IsSetFlag(n *int32) bool {
	v := atomic.LoadInt32(n)
	if v == setFlag {
		return true
	}

	return false
}

func IsUnsetFlag(n *int32) bool {
	v := atomic.LoadInt32(n)
	if v == setFlag {
		return false
	}

	return true
}
