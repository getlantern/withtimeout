// package withtimeout provides functionality for performing operations with
// a timeout.
package withtimeout

import (
	"time"
)

const (
	timeoutErrorString = "igdman: Operation timed out"
)

type timeoutError struct{}

func (timeoutError) Error() string { return timeoutErrorString }

// Do executes the given fn and returns either the result of executing it or an
// error if fn did not complete within timeout. If execution timed out, timedOut
// will be true.
func Do(timeout time.Duration, fn func() (interface{}, error)) (result interface{}, timedOut bool, err error) {
	return DoOr(timeout, fn, nil)
}

// DoOr is like Do but also executes the given onTimeout function if and when fn
// times out.
func DoOr(timeout time.Duration, fn func() (interface{}, error), onTimeout func()) (result interface{}, timedOut bool, err error) {
	resultCh := make(chan *resultWithError)

	go func() {
		result, err := fn()
		select {
		case resultCh <- &resultWithError{result, err}:
			// result submitted
		}
	}()

	select {
	case <-time.After(timeout):
		if onTimeout != nil {
			onTimeout()
		}
		return nil, true, timeoutError{}
	case rwe := <-resultCh:
		return rwe.result, false, rwe.err
	}
}

type resultWithError struct {
	result interface{}
	err    error
}
