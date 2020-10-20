package goterminal

import (
	"testing"
	"time"
)

func TestPrintPercent(t *testing.T) {
	var bar Bar
	bar.NewOption(0, 50)
	for i := 0; i <= 100; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Show(int64(i))
	}
	bar.Finish()
}
