package goid

import (
	"bytes"
	"runtime"
	"strconv"
)

var (
	openFlag = false
)

func SetOpenGID() {
	openFlag = true
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
