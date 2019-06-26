package main

import (
	"fmt"

	"github.com/jan25/gocui"
)

// Stats is to keep track of
// stats view
type Stats struct {
	// name of the view
	name string

	// position and dims
	x, y int
	w, h int

	// timer instance
	timer *Timer
}

func newStatsView(name string, x, y int, w, h int) *Stats {
	return &Stats{
		name:  name,
		x:     x,
		y:     y,
		w:     w,
		h:     h,
		timer: NewTimer(),
	}
}

// Layout manager for stats view
func (s *Stats) Layout(g *gocui.Gui) error {
	v, err := g.SetView(s.name, s.x, s.y, s.x+s.w, s.y+s.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if s.timer.IsActive() {
		s.updateTime(v)
	}
	return nil
}

func (s *Stats) updateTime(v *gocui.View) error {
	elapsedTime, err := s.timer.ElapsedTime()
	if err != nil {
		return err
	}

	v.Clear()
	fmt.Fprintf(v, "%02d:%02d", elapsedTime.Mins, elapsedTime.Secs)
	return nil
}

// StartTimer starts the timer
func (s *Stats) StartTimer() error {
	err := s.timer.Start()
	return err
}

// StopTimer stops the timer
func (s *Stats) StopTimer() error {
	err := s.timer.Stop()
	return err
}
