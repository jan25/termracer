package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/jan25/gocui"
)

// OneStat represents information
// about completed race
type OneStat struct {
	// words per minute
	Wpm int
	// accuracy eg. 95.6%
	Accuracy float64
	// time of race
	When time.Time
}

// Stats is a datastructure to
// store stats for past races
type Stats struct {
	History  []*OneStat
	Selected int // To keep track of highlighted stat in history

	Current *OneStat
}

// InitNewStat initializes current OneStat
// used for beginning a new race
func (s *Stats) InitNewStat() {
	if s.Current == nil {
		s.Current = &OneStat{
			Wpm:      0,
			Accuracy: float64(100.00),
			When:     time.Now(),
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
	// Always selected most recent race stat on finish
	s.Selected = len(s.History) - 1
	return nil
}

// ScrollDown is a keybinding
// Increment selected stat index
func (s *Stats) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	if s.Selected+1 < len(s.History) {
		s.Selected++
	} else {
		Logger.Info("End of history reached. Can not scroll further up")
	}
	// just to adhere to KeyBinding handler interface
	return nil
}

// ScrollUp is a keybinding
// Decrements selected stat index
func (s *Stats) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	if s.Selected > 0 {
		s.Selected--
	} else {
		Logger.Info("Top of history reached. Can not scroll further down")
	}
	// just to adhere to KeyBinding handler interface
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
		v.Title = "Race in progress"
		s.updateRaceStats(v)
	} else {
		v.Title = "Race History"
		g.SetCurrentView(s.name)
		s.showRecentRaceStats(v)
	}

	return nil
}

func (s *StatsView) showRecentRaceStats(v *gocui.View) {
	v.Clear()

	if len(s.stats.History) == 0 {
		v.Wrap = true
		fmt.Fprintf(v, "No races available to show")
		return
	}

	fmt.Fprintln(v, "\033[0;4mwpm acc. when    \033[0m")

	selected := s.stats.Selected
	for i := selected; i >= 0; i-- {
		stat := s.stats.History[i]
		f := "%s\n"
		if i == s.stats.Selected {
			f = "\033[0;7m%s\033[0m\n"
		}
		fmt.Fprintf(v, f,
			fmt.Sprintf("%-3d %3d%% %-8s", stat.Wpm, int(stat.Accuracy), FormatDate(stat.When)))
	}
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
	a, err := CalculateAccuracy(paragraph.CharsUptoCurrent(), word.Mistyped)
	// Only update accuracy if we were able to calculate it
	// This is to handle 0 chars types at start of race
	if err == nil {
		s.stats.Current.Accuracy = a
	}

	currentStat := s.stats.Current

	v.Clear()
	fmt.Fprintf(v, "%02d:%02d \n", elapsedTime.Mins, elapsedTime.Secs)
	if currentStat != nil {
		fmt.Fprintf(v, "wpm %d \nAccuracy %.2f%% \n", currentStat.Wpm, currentStat.Accuracy)
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
// finished indicates if race is finished but not ended
func (s *StatsView) StopRace(finished bool) error {
	if finished {
		s.stats.FinishCurrent()
	}
	err := s.timer.Stop()
	return err
}

// InitKeyBindings adds keybindings to scroll in stats history
func (s *StatsView) InitKeyBindings(g *gocui.Gui) {
	if err := g.SetKeybinding(s.name, gocui.KeyCtrlJ, gocui.ModNone, s.stats.ScrollUp); err != nil {
		Logger.Warn(fmt.Sprintf("%v", err))
	}
	if err := g.SetKeybinding(s.name, gocui.KeyCtrlK, gocui.ModNone, s.stats.ScrollDown); err != nil {
		Logger.Warn(fmt.Sprintf("%v", err))
	}
}
