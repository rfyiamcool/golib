package main

import (
	"fmt"

	"github.com/rfyiamcool/golib/ringqueue"
)

func main() {
	q := ringqueue.New()

	for i := 0; i < 21; i++ {
		q.Add(i)
	}

	for i := 0; i < 21; i++ {
		if q.Peek().(int) != i {
			fmt.Println(q.Peek())
		}

		x := q.Remove()
		if x != i {
			fmt.Println(x)
		}
	}
}
