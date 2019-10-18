package main

import (
	"errors"
	"sync"
	"time"
)

// Timer is a stopwatch like functionality
// Doesn't use ticker, but keeps tracks of
// start time to figure the elapsed time
type Timer struct {
	start time.Time

	// Keeps track of state of timer
	active bool
	done   chan struct{}
	wg     sync.WaitGroup
}

// TimeFormatted wraps time.Duration converted to
// mins:seconds format
type TimeFormatted struct {
	Mins, Secs int
}

// NewTimer creates and returns new timer instance
func NewTimer() *Timer {
	return &Timer{}
}

// Start starts the timer
func (t *Timer) Start() error {
	if t.active {
		return errors.New("timer already started")
	}

	t.start = time.Now()
	t.active = true
	t.done = make(chan struct{})

	t.wg.Add(1)
	return nil
}

func (t *Timer) getDoneCh() chan struct{} {
	if t.done == nil {
		t.done = make(chan struct{})
	}
	return t.done
}

// Stop stops the timer
func (t *Timer) Stop() error {
	t.start = time.Now()
	t.active = false

	select {
	case <-t.getDoneCh():
		// channel already closed
		return errors.New("timer already stopped")
	default:
		close(t.done)
	}

	return nil
}

// IsActive returns true if timer is active
func (t *Timer) IsActive() bool {
	return t.active
}

// ElapsedTime is time the timer has been active for
// returns error if timer is not active
func (t *Timer) ElapsedTime() (*TimeFormatted, error) {
	elapsed, err := t.elapsedDuration()
	if err != nil {
		return nil, err
	}

	ms := elapsed / time.Millisecond
	secs := ms / 1000
	mins := secs / 60

	secs = secs % 60

	tf := TimeFormatted{
		Mins: int(mins),
		Secs: int(secs),
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
