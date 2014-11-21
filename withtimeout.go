// package withtimeout provides functionality for performing operations with
// a timeout.
package withtimeout

import (
	"time"
)

const (
	timeoutErrorString = "withtimeout: Operation timed out"
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
func DoOr(timeout time.Duration, fn func() (interface{}, error), onTimeout func() error) (result interface{}, timedOut bool, err error) {
	resultCh := make(chan *resultWithError, 1)

	go func() {
		result, err := fn()
		resultCh <- &resultWithError{result, err}
	}()

	select {
	case <-time.After(timeout):
		var err error = timeoutError{}
		if onTimeout != nil {
			err = onTimeout()
		}
		return nil, true, err
	case rwe := <-resultCh:
		return rwe.result, false, rwe.err
	}
}

type resultWithError struct {
	result interface{}
	err    error
}
