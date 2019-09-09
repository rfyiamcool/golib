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

func Add(delay time.Duration, callback func()) (*Task, error) {
	return DefaultTimeWheel.Add(delay, callback)
}

func AddCron(delay time.Duration, callback func()) (*Task, error) {
	return DefaultTimeWheel.AddCron(delay, callback)
}

func Remove(task *Task) error {
	return DefaultTimeWheel.Remove(task)
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
