package current

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type BatchHandler[T, R any] func(ctx context.Context, req T) (R, error)

type Options struct {
	// max size of current worker
	syncLimit int64
	// max time to wait for
	timeout time.Duration
	// request rate limit
	rateLimit *rate.Limiter
	// if skip unfinish param when handler returns error.
	// default as false, will wait for all requests done or timeout
	fastFails bool
	// when request reach limit, should we just record error or waiting fo token.
	// default as false, will wait for tokens to exec
	rateWithError bool
	// handlers for returns error translate or collected
	errHandler func(err error) error
}
type Option func(*Options)

// WithTimeout max time to wait for
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) { o.timeout = timeout }
}

// WithFastFails if skip params haven't been handle when handler returns error.
// default as false, will wait for all requests done or timeout
func WithFastFails(fastFails bool) Option {
	return func(o *Options) { o.fastFails = fastFails }
}

// WithReturnErrHandler handlers for returns error translate or collected
func WithReturnErrHandler(errHandler func(err error) error) Option {
	return func(o *Options) { o.errHandler = errHandler }
}

// WithRateLimit request rate limit
func WithRateLimit(qpsLimit float64, boostLimit int) Option {
	return func(o *Options) {
		o.rateLimit = rate.NewLimiter(rate.Limit(qpsLimit), boostLimit)
	}
}

// WithSyncLimit max size of current worker
func WithSyncLimit(syncLimit int64) Option {
	return func(o *Options) {
		if syncLimit > 0 {
			o.syncLimit = syncLimit
		}
	}
}

// WithRateWithError when request reach limit, should we just record error or waiting fo token.
// default as false, will wait for tokens to exec
func WithRateWithError(rateWithError bool) Option {
	return func(o *Options) { o.rateWithError = rateWithError }
}
