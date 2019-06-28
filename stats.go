package main

import (
	"errors"
	"fmt"

	"text/tabwriter"

	"github.com/jan25/gocui"
)

// OneStat represents information
// about completed race
type OneStat struct {
	// words per minute
	Wpm int
	// accuracy eg. 95.6%
	Accuracy float64
}

// Stats is datastructure to
// store stats for past races
type Stats struct {
	History []*OneStat
	Current *OneStat
}

// InitNewStat initializes current OneStat
// used for beginning a new race
func (s *Stats) InitNewStat() {
	if s.Current == nil {
		s.Current = &OneStat{
			Wpm:      0,
			Accuracy: float64(100.00),
		}
	}
	if s.History == nil {
		s.History = make([]*OneStat, 0)
	}
}

// FinishCurrent adds current Stat
// to history closing the current race
func (s *Stats) FinishCurrent() error {
	if s.Current == nil {
		return errors.New("No current Stat to finish")
	}
	s.History = append(s.History, s.Current)
	s.Current = nil
	return nil
}

// StatsView is to keep track of
// stats and the view
type StatsView struct {
	// name of the view
	name string

	// position and dims
	x, y int
	w, h int

	// timer instance
	timer *Timer

	// stats instance
	stats *Stats
}

func newStatsView(name string, x, y int, w, h int) *StatsView {
	return &StatsView{
		name:  name,
		x:     x,
		y:     y,
		w:     w,
		h:     h,
		timer: NewTimer(),
		stats: &Stats{},
	}
}

// Layout manager for stats view
func (s *StatsView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(s.name, s.x, s.y, s.x+s.w, s.y+s.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if s.timer.IsActive() {
		s.updateRaceStats(v)
	} else {
		s.showRecentRaceStats(v)
	}

	return nil
}

func (s *StatsView) showRecentRaceStats(v *gocui.View) {
	v.Clear()

	w := new(tabwriter.Writer)
	w.Init(v, 0, 9, 0, '\t', 0)
	fmt.Fprintln(w, "No. WPM ACCURACY")
	for i := len(s.stats.History) - 1; i >= 0; i-- {
		stat := s.stats.History[i]
		fmt.Fprintln(w,
			fmt.Sprintf("%d %d %.2f%%", (i+1), stat.Wpm, stat.Accuracy))
	}
	w.Flush()
}

func (s *StatsView) updateRaceStats(v *gocui.View) error {
	elapsedTime, err := s.timer.ElapsedTime()
	if err != nil {
		return err
	}

	secs := elapsedTime.Mins*60 + elapsedTime.Secs
	if secs != 0 {
		s.stats.Current.Wpm = CalculateWpm(paragraph.CountDoneWords(), secs)
	} else {
		s.stats.Current.Wpm = 0
	}
	s.stats.Current.Accuracy = CalculateAccuracy(paragraph.CharsUptoCurrent(), word.Mistyped)

	currentStat := s.stats.Current

	v.Clear()
	fmt.Fprintf(v, "%02d:%02d \n", elapsedTime.Mins, elapsedTime.Secs)
	if currentStat != nil {
		fmt.Fprintf(v, "WPM %d \nACCURACY %.2f%% \n", currentStat.Wpm, currentStat.Accuracy)
	}
	return nil
}

// StartRace starts the timer
func (s *StatsView) StartRace() error {
	err := s.timer.Start()
	s.stats.InitNewStat()
	return err
}

// StopRace stops the timer
func (s *StatsView) StopRace() error {
	err := s.timer.Stop()
	s.stats.FinishCurrent()
	return err
}
