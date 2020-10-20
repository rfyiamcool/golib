package main

import (
	"time"

	"github.com/rfyiamcool/golib/goterminal"
)

func main() {
	var bar goterminal.Bar
	bar.NewOption(0, 1000)
	for i := 0; i <= 100; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Show(int64(i))
	}
	bar.Finish()
}
