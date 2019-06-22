package main

import (
	"errors"
	"time"
)

// Timer is a stopwatch like functionality
// Doesn't use ticker, but keeps tracks of
// start time to figure the elapsed time
type Timer struct {
	start  time.Time
	active bool
}

// TimeFormatted wraps time.Duration converted to
// seconds:centiseconds format
type TimeFormatted struct {
	Seconds, CentiSeconds int
}

// NewTimer creates and returns new timer instance
func NewTimer() *Timer {
	var timer Timer
	return &timer
}

// Start starts the timer
func (t *Timer) Start() {
	t.start = time.Now()
	t.active = true
}

// Stop stops the timer
func (t *Timer) Stop() {
	t.start = time.Now()
	t.active = false
}

// ElapsedTime is time the timer is active
// returns error if time is not active
func (t *Timer) ElapsedTime() (*TimeFormatted, error) {
	elapsed, err := t.elapsedDuration()
	if err != nil {
		return nil, err
	}

	ms := elapsed / time.Millisecond
	secs := ms / 1000
	csecs := (ms / 10) % 100

	tf := TimeFormatted{
		Seconds:      int(secs),
		CentiSeconds: int(csecs),
	}
	return &tf, nil
}

func (t *Timer) elapsedDuration() (time.Duration, error) {
	if t.active {
		return time.Since(t.start), nil
	}

	// return ~0 seconds
	// is there better way to return `nil` time?
	return time.Microsecond, errors.New("timer is not active")
}
