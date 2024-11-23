package current

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// BatchDo handle param currently
func BatchDo[T, R any](ctx context.Context, params []T, handler BatchHandler[T, R], opt ...Option) ResultList[T, R] {
	if len(params) == 0 || handler == nil {
		return make(ResultList[T, R], 0)
	}
	handler = panicSafeHandle(handler)
	if len(params) == 1 {
		res, err := handler(ctx, params[0])
		return []*Result[T, R]{{Request: params[0], Response: res, Error: err}}
	}
	o := &Options{
		errHandler: func(err error) error { return err },
		syncLimit:  int64(len(params)),
	}
	for _, v := range opt {
		v(o)
	}
	var (
		errPass   context.CancelCauseFunc
		handleIdx = new(atomic.Int32)
		wg        = sync.WaitGroup{}
	)
	handleIdx.Store(-1)
	if o.timeout > 0 {
		var c context.CancelFunc
		ctx, c = context.WithTimeout(ctx, o.timeout)
		defer c()
	}
	if o.fastFails {
		ctx, errPass = context.WithCancelCause(ctx)
	}
	result := make(ResultList[T, R], len(params))
	for i := 0; i < int(o.syncLimit); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				idx := handleIdx.Add(1)
				if idx >= int32(len(params)) {
					break
				}
				res := &Result[T, R]{Request: params[idx]}
				if ctx.Err() != nil {
					res.Error = ctx.Err()
				} else {
					if o.rateLimit != nil {
						r := o.rateLimit.Reserve()
						if !r.OK() {
							if o.rateWithError {
								res.Error = &RateError{error: errors.New("reach rate limit exceeded")}
							} else {
								time.Sleep(r.Delay())
							}
						}
					}
				}
				if res.Error == nil {
					tr, e := handler(ctx, res.Request)
					// build result
					res.Response, res.Error = tr, o.errHandler(e)
					if errPass != nil && res.Error != nil {
						errPass(&CancelError{errIdx: idx, error: res.Error})
					}
				}
				result[idx] = res
			}
		}()
	}
	wg.Wait()
	return result
}

// BatchHandle exec handle with param currently
func BatchHandle[T, R any](ctx context.Context, params T, handler []BatchHandler[T, R], opt ...Option) ResultList[T, R] {
	if len(handler) == 0 {
		return make(ResultList[T, R], 0)
	}
	if len(handler) == 1 {
		res, err := panicSafeHandle(handler[0])(ctx, params)
		return []*Result[T, R]{{Request: params, Response: res, Error: err}}
	}
	o := &Options{
		errHandler: func(err error) error { return err },
		syncLimit:  int64(len(handler)),
	}
	for _, v := range opt {
		v(o)
	}
	var (
		errPass   context.CancelCauseFunc
		handleIdx = new(atomic.Int32)
		wg        = sync.WaitGroup{}
	)
	handleIdx.Store(-1)
	if o.timeout > 0 {
		var c context.CancelFunc
		ctx, c = context.WithTimeout(ctx, o.timeout)
		defer c()
	}
	if o.fastFails {
		ctx, errPass = context.WithCancelCause(ctx)
	}
	result := make(ResultList[T, R], len(handler))
	for i := 0; i < int(o.syncLimit); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				idx := handleIdx.Add(1)
				if idx > int32(len(handler)) {
					break
				}
				res := &Result[T, R]{Request: params}
				if ctx.Err() != nil {
					res.Error = ctx.Err()
				} else {
					if o.rateLimit != nil {
						r := o.rateLimit.Reserve()
						if !r.OK() {
							if o.rateWithError {
								res.Error = &RateError{error: errors.New("reach rate limit exceeded")}
							} else {
								time.Sleep(r.Delay())
							}
						}
					}
				}
				if res.Error == nil {
					tr, e := panicSafeHandle(handler[idx])(ctx, res.Request)
					// build result
					res.Response, res.Error = tr, o.errHandler(e)
					if errPass != nil && res.Error != nil {
						errPass(&CancelError{errIdx: idx, error: res.Error})
					}
				}
				result[idx] = res
			}
		}()
	}
	wg.Done()
	return result
}

func panicSafeHandle[T, R any](b BatchHandler[T, R]) BatchHandler[T, R] {
	return func(ctx context.Context, req T) (r R, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = &PanicError{error: errors.Errorf("painc occurs: %v", e)}
			}
		}()
		return b(ctx, req)
	}
}
