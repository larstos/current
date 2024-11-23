package current

import "fmt"

type PanicError struct{ error }
type RateError struct{ error }
type CancelError struct {
	error
	errIdx int32
}

func (e *CancelError) Error() string {
	return fmt.Sprintf("task failed idx[%d]:%s", e.errIdx, e.error.Error())
}
func (e *CancelError) GetErrIdx() int32 { return e.errIdx }
func (e *CancelError) GetError() error  { return e.error }

type Result[T, R any] struct {
	Request  T
	Response R
	Error    error
}

type ResultList[T, R any] []*Result[T, R]

func (r ResultList[T, R]) Len() int { return len(r) }
func (r ResultList[T, R]) ExportsRes() []R {
	res := make([]R, 0, len(r))
	for _, v := range r {
		res = append(res, v.Response)
	}
	return res
}

func (r ResultList[T, R]) FindByRequest(request T, equal func(a, b T) bool) *Result[T, R] {
	if equal == nil {
		return nil
	}
	for _, v := range r {
		if equal(v.Request, request) {
			return v
		}
	}
	return nil
}
