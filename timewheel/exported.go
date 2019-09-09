package timewheel

import (
	"time"
)

var (
	DefaultTimeWheel, _ = NewTimeWheel(time.Second, 120)
)

func ResetDefaultTimeWheel(tw *TimeWheel) {
	DefaultTimeWheel = tw
}

func NewTimer(delay time.Duration) *Timer {
	return DefaultTimeWheel.NewTimer(delay)
}

func NewTicker(delay time.Duration) *Ticker {
	return DefaultTimeWheel.NewTicker(delay)
}

func After(delay time.Duration) {
	DefaultTimeWheel.After(delay)
}

func Sleep(delay time.Duration) {
	DefaultTimeWheel.Sleep(delay)
}
