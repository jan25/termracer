package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	timer := NewTimer()
	timer.Start()

	time.Sleep(1200 * time.Millisecond)
	elapsed, err := timer.ElapsedTime()
	if err != nil {
		t.Fatal("error getting elapsed time")
	}

	assert.True(t, elapsed.Secs == 1,
		fmt.Sprintf("timer has incorrect elapsed secs %v", elapsed))

	timer.Stop()
	elapsed, err = timer.ElapsedTime()
	assert.Error(t, err, "timer must throw error if not active")
}
