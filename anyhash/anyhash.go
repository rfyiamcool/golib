package anyhash

import (
	"unsafe"

	"github.com/mitchellh/hashstructure"
)

func Hash(key interface{}) uint64 {
	switch k := key.(type) {
	case int:
		return uint64(k)
	case int32:
		return uint64(k)
	case uint32:
		return uint64(k)
	case int64:
		return uint64(k)
	case uint:
		return uint64(k)
	case uint64:
		return k
	case string:
		return memHashString(k)
	case []byte:
		return memHash(k)
	case byte:
		return memHash([]byte{k})
	default:
		code, err := hashstructure.Hash(key, nil)
		if err != nil {
			panic(err)
		}
		return code
	}
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

func memHash(data []byte) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&data))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}

func memHashString(str string) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}
