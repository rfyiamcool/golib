package where

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Printf("line:9 ----> %s\n", Get())
	fmt.Printf("%d\n", mock(4))
}

func mock(num int) int {
	fmt.Printf("line:14 ---> %s\n", Get())
	if num == 1 {
		return 1
	} else {
		return num * mock(num-1)
	}
}
