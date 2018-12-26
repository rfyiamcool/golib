package retry

import (
	"context"
	"time"
)

// RetriableErr is an error type which can be retried
type RetriableErr struct {
	error
}

// Retriable makes an error be retriable
func Retriable(err error) *RetriableErr {
	return &RetriableErr{err}
}

// Retry ensures that the do function will be executed until some condition being satisfied
type Retry struct {
	base time.Duration
}

var r = New()

func (r *Retry) ensure(ctx context.Context, times int, do func() error) error {
	var (
		duration  = r.base
		alwayFlag = true
		err       error
	)
	if times > 0 {
		alwayFlag = false
	}

	for {
		if !alwayFlag && times == 0 {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := do(); err != nil {
			if !alwayFlag {
				times--
			}
			if _, ok := err.(*RetriableErr); ok {
				if r.base > 0 {
					time.Sleep(duration)
				}
				continue
			}
			return err
		}

		return nil
	}
}

// Ensure keeps retring until ctx is done
func (r *Retry) Ensure(ctx context.Context, do func() error) error {
	return r.ensure(ctx, 0, do)
}

func (r *Retry) EnsureRetryTimes(ctx context.Context, times int, do func() error) error {
	return r.ensure(ctx, times, do)
}

// Option is an option to new a Retry object
type Option func(r *Retry)

// BackoffStrategy defines the backoff strategy of retry
type BackoffStrategy func(last time.Duration) time.Duration

// WithBaseDelay set the first delay duration, default 10ms
func WithBaseDelay(base time.Duration) Option {
	return func(r *Retry) {
		r.base = base
	}
}

// New a retry object
func New(opts ...Option) *Retry {
	r := &Retry{base: 10 * time.Millisecond}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Ensure keeps retring until ctx is done, it use a default retry object
func Ensure(ctx context.Context, do func() error) error {
	return r.Ensure(ctx, do)
}
