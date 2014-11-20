package withtimeout

import (
	"fmt"
	"testing"
	"time"

	"github.com/getlantern/testify/assert"
)

var (
	expectedText = "My Text"
	expectedErr  = fmt.Errorf("My Error")
)

func TestSuccess(t *testing.T) {
	text, timedOut, err := Do(1*time.Second, func() (interface{}, error) {
		return expectedText, expectedErr
	})
	assert.False(t, timedOut, "Should not have timed out")
	assert.Equal(t, expectedText, text, "Text should match expected")
	assert.Equal(t, expectedErr, err, "Error should match expected")
}

func TestTimeout(t *testing.T) {
	calledBack := false
	text, timedOut, err := DoOr(10*time.Millisecond, func() (interface{}, error) {
		time.Sleep(11 * time.Millisecond)
		return expectedText, expectedErr
	}, func() {
		calledBack = true
	})
	assert.True(t, timedOut, "Should have timed out")
	assert.True(t, calledBack, "Should have called back after timing out")
	assert.NotNil(t, err, "There should be an error")
	assert.Nil(t, text, "Text should be nil")
	assert.Equal(t, timeoutErrorString, err.Error(), "Error should contain correct string")
}
