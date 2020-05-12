package errgroup

import (
	"context"
	"sync"
)

// modify err returned
// refer golang.org/x/sync/errgroup

type Group struct {
	cancel func()

	wg sync.WaitGroup

	sync.Mutex
	errs []error
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
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
