package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	defaultCtx, defaultCancel = context.WithCancel(context.Background())
	defaultDelayInterval      = time.Duration(100 * time.Millisecond)
	defaultRetriableError     = RetriableMesg("retriable")
)

// RetriableErr is an error type which can be retried
type RetriableErr struct {
	error
}

// Retriable makes an error be retriable
func Retriable(err error) *RetriableErr {
	return &RetriableErr{err}
}

func RetriableMesg(mesg string) *RetriableErr {
	return &RetriableErr{errors.New(mesg)}
}

// Retry ensures that the do function will be executed until some condition being satisfied
type Retry struct {
	ctx      context.Context
	base     time.Duration
	backoff  *Backoff // if backoff not nil, use backoff, ignore base duration
	recovery bool
}

func (r *Retry) ensure(times int, do func() error) error {
	var (
		alway = true
		err   error
	)
	if times > 0 {
		alway = false
	}

	for {
		if !alway && times == 0 {
			return err
		}
		if r.isExited() {
			return r.ctx.Err()
		}

		err = r.handle(do)
		if err == nil {
			return nil
		}

		if !alway {
			times--
		}
		if _, ok := err.(*RetriableErr); ok {
			r.sleep()
			continue
		}
		return err
	}
}

func (r *Retry) handle(fn func() error) (err error) {
	if !r.recovery {
		return fn()
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovery %v", r)
		}
	}()

	return fn()
}

func (r *Retry) isExited() bool {
	select {
	case <-r.ctx.Done():
		return true
	default:
		return false
	}
}

func (r *Retry) sleep() {
	var duratime = r.base

	if r.backoff != nil {
		duratime = r.backoff.Duration()
	}

	select {
	case <-time.After(duratime):
	case <-r.ctx.Done():
	}
}

// Ensure keeps retring until ctx is done
func (r *Retry) Ensure(do func() error) error {
	return r.ensure(0, do)
}

// retry times limit
func (r *Retry) EnsureRetryTimes(times int, do func() error) error {
	return r.ensure(times, do)
}

// Option is an option to new a Retry object
type Option func(r *Retry)

// WithBaseDelay set the first delay duration, default 10ms
func WithBaseDelay(base time.Duration) Option {
	return func(r *Retry) {
		r.base = base
	}
}

func WithCtx(ctx context.Context) Option {
	return func(r *Retry) {
		r.ctx = ctx
	}
}

func WithRecovery() Option {
	return func(r *Retry) {
		r.recovery = true
	}
}

func WithBackoff(bo *Backoff) Option {
	return func(r *Retry) {
		r.backoff = bo
	}
}

type Backoff struct {
	MinDelay time.Duration
	MaxDelay time.Duration
	Factor   float64
	Jitter   bool
	attempts float64
}

func (b *Backoff) Duration() time.Duration {
	dur := float64(b.MinDelay) * math.Pow(b.Factor, b.attempts)
	if b.Jitter == true {
		dur = rand.Float64()*(dur-float64(b.MinDelay)) + float64(b.MinDelay)
	}
	if dur > float64(b.MaxDelay) {
		return b.MaxDelay
	}

	b.attempts++
	return time.Duration(dur)
}

// New a retry object
func New(opts ...Option) *Retry {
	r := &Retry{
		ctx:  defaultCtx,
		base: defaultDelayInterval,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Ensure keeps retring until ctx is done, it use a default retry object
func Ensure(ctx context.Context, do func() error) error {
	r := &Retry{
		ctx:  ctx,
		base: defaultDelayInterval,
	}
	return r.Ensure(do)
}

// ensure backoff
func EnsureWithBackoff(ctx context.Context, bo *Backoff, do func() error) error {
	r := &Retry{
		ctx:     ctx,
		base:    defaultDelayInterval,
		backoff: bo,
	}
	return r.Ensure(do)
}
