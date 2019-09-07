package timewheel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type A struct {
	a int
	b string
}

func callback() {
	fmt.Println("callback !!!")
}

func newTimeWheel() *TimeWheel {
	tw, err := NewTimeWheel(1*time.Second, 360)
	if err != nil {
		panic(err)
	}
	tw.Start()
	return tw
}

func TestAdd(t *testing.T) {
	tw := newTimeWheel()
	_, err := tw.Add(time.Second*1, callback)
	if err != nil {
		t.Fatalf("test add failed, %v", err)
	}
	time.Sleep(time.Second * 5)
	tw.Stop()
}

func TestCron(t *testing.T) {
	tw := newTimeWheel()
	_, err := tw.AddCron(time.Second*1, callback)
	if err != nil {
		t.Fatalf("test add failed, %v", err)
	}
	time.Sleep(time.Second * 5)
	tw.Stop()
}

func TestTicker(t *testing.T) {
	tw := newTimeWheel()
	ticker := tw.NewTicker(time.Second * 1)
	go func() {
		time.Sleep(5 * time.Second)
		ticker.Stop()
		fmt.Println("call stop")
	}()
	for {
		select {
		case <-ticker.C:
			callback()
		case <-ticker.Ctx.Done():
			return
		}
	}
}

func TestBatchTicker(t *testing.T) {
	tw := newTimeWheel()
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := tw.NewTicker(time.Second * 1)
			go func() {
				time.Sleep(5 * time.Second)
				ticker.Stop()
				fmt.Println("call stop")
			}()
			for {
				select {
				case <-ticker.C:
					callback()
				case <-ticker.Ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestTimerReset(t *testing.T) {
	tw := newTimeWheel()
	timer := tw.NewTimer(1 * time.Second)
	now := time.Now()
	<-timer.C
	t.Logf(time.Since(now).String())

	timer.Reset(2 * time.Second)
	now = time.Now()
	<-timer.C
	t.Logf(time.Since(now).String())

	now = time.Now()
	timer.Reset(3 * time.Second)
	<-timer.C
	t.Logf(time.Since(now).String())

	now = time.Now()
	timer.Reset(5 * time.Second)
	<-timer.C
	t.Logf(time.Since(now).String())
}

func TestRemove(t *testing.T) {
	tw := newTimeWheel()
	task, err := tw.Add(time.Second*1, callback)
	if err != nil {
		t.Fatalf("test add failed, %v", err)
	}
	tw.Remove(task)
	time.Sleep(time.Second * 5)
	tw.Stop()
}

func TestHwTimer(t *testing.T) {
	tw := newTimeWheel()
	worker := 10
	delay := 5

	wg := sync.WaitGroup{}
	for index := 0; index < worker; index++ {
		wg.Add(1)
		var (
			htimer = tw.NewTimer(time.Duration(delay) * time.Second)
			maxnum = 20
			incr   = 0
		)
		go func(idx int) {
			defer wg.Done()
			for incr < maxnum {
				now := time.Now()
				target := time.Now().Add(time.Duration(delay) * time.Second)
				select {
				case <-htimer.C:
					htimer.Reset(time.Duration(delay) * time.Second)
					end := time.Now()
					if end.Before(target.Add(-1 * time.Second)) {
						t.Log("before 1s run")
					}
					if end.After(target.Add(1 * time.Second)) {
						t.Log("delay 1s run")
					}
					fmt.Println("id: ", idx, "cost: ", time.Now().Sub(now))
				}
				incr++
			}
		}(index)
	}
	wg.Wait()
}

func BenchmarkAdd(b *testing.B) {
	tw := newTimeWheel()
	for i := 0; i < b.N; i++ {
		_, err := tw.Add(time.Second, callback)
		if err != nil {
			b.Fatalf("benchmark Add failed, %v", err)
		}
	}
}
