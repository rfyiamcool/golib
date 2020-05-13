package errgroup

import (
	"context"
	"sync"
)

// modify err returned
// refer golang.org/x/sync/errgroup

func ConvertQueue(l interface{}) chan interface{} {
	s := reflect.ValueOf(l)
	c := make(chan interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		c <- s.Index(i).Interface()
	}
	close(c)
	return c
}

type Group struct {
	cancel func()

	wg sync.WaitGroup

	sync.Mutex
	errs []error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) WaitOK() ([]error, bool) {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.errs, len(g.errs) == 0
}

func (g *Group) Spawn(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.Lock()
			g.errs = append(g.errs, err)
			if g.cancel != nil {
				g.cancel()
			}
			g.Unlock()
		}
	}()
}
