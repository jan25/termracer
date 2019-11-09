package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jan25/color"
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/pkg/utils"
	"go.uber.org/zap"
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

// LoadHistory loads race history from local
// file storage into inmemory History array
func (s *Stats) LoadHistory() error {
	fname, _ := GetHistoryFilePath()
	f, err := os.Open(fname)
	if err != nil {
		Logger.Error("failed to open racehistory.csv file")
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)

	// column name to index mapper
	cmapper := make(map[string]int)
	records, _ := r.ReadAll()
	if len(records) == 0 {
		// Append column names as first line
		// to handle the case of newly created file
		utils.AppendLineEOF(fname, "wpm,acc,when")
	}
	for i, record := range records {
		if i == 0 {
			// Read column names
			for j, cname := range record {
				cmapper[cname] = j
			}
		} else {
			if s.History == nil {
				s.History = make([]*OneStat, 0)
			}
			wpm, err := strconv.Atoi(record[cmapper["wpm"]])
			if err != nil {
				Logger.Warn("failed to read a record: " + strings.Join(record, ""))
				continue
			}
			acc, err := strconv.ParseFloat(record[cmapper["acc"]], 64)
			if err != nil {
				Logger.Warn("failed to read a record: " + strings.Join(record, ""))
				continue
			}
			when, err := time.Parse("02/01/06", record[cmapper["when"]])

			stat := &OneStat{
				Wpm:      wpm,
				Accuracy: acc,
				When:     when,
			}
			s.History = append(s.History, stat)
		}
	}
	Logger.Info("Finished loading records from file", zap.Int("records", len(records)))
	s.Selected = len(s.History) - 1

	return nil
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
	s.AppendToFile(s.Current) // append to storage file
	s.Current = nil
	// Always selected most recent race stat on finish
	s.Selected = len(s.History) - 1
	return nil
}

// AppendToFile appends last finished race to
// localstorage file
func (s *Stats) AppendToFile(stat *OneStat) error {
	line := fmt.Sprintf("%d,%f,%s", stat.Wpm, stat.Accuracy, utils.FormatDate(stat.When))
	f, _ := GetHistoryFilePath()
	if err := utils.AppendLineEOF(f, line); err != nil {
		Logger.Warn("Failed to append to file", zap.Error(err))
		return err
	}
	Logger.Info("Successfuly appended to file")
	return nil
}

// ScrollDown is a keybinding
// increments selected stat index
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
// decrements selected stat index
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
	stats := &Stats{}
	stats.LoadHistory()
	return &StatsView{
		name:  name,
		x:     x,
		y:     y,
		w:     w,
		h:     h,
		timer: NewTimer(),
		stats: stats,
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
		cyan := color.New(color.FgCyan)
		cyan.Fprintf(v, "No races available to show")
		return
	}

	// Heading
	ul := color.New(color.Underline)
	ul.Fprintln(v, "wpm acc. when    ")

	selected := s.stats.Selected
	for i := selected; i >= 0; i-- {
		stat := s.stats.History[i]
		f := "%s\n"
		if i == s.stats.Selected {
			white := color.New(color.BgWhite)
			white.Fprintf(v, f,
				fmt.Sprintf("%-3d %3d%% %-8s", stat.Wpm, int(stat.Accuracy), utils.FormatDate(stat.When)))
		} else {
			fmt.Fprintf(v, f,
				fmt.Sprintf("%-3d %3d%% %-8s", stat.Wpm, int(stat.Accuracy), utils.FormatDate(stat.When)))
		}
	}
}

func (s *StatsView) updateRaceStats(v *gocui.View) error {
	elapsedTime, err := s.timer.ElapsedTime()
	if err != nil {
		return err
	}

	secs := elapsedTime.Mins*60 + elapsedTime.Secs
	if secs != 0 {
		s.stats.Current.Wpm = utils.CalculateWpm(paragraph.CountDoneWords(), secs)
	} else {
		s.stats.Current.Wpm = 0
	}
	a, err := utils.CalculateAccuracy(paragraph.CharsUptoCurrent(), word.Mistyped)
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
	if err != nil {
		// this is to fix pressing ctrlS multiple times
		// which caused crashes when ending with ctrlE
		// due to how timer.wg works
		return err
	}
	go s.updateTimer(s.timer, g)
	s.stats.InitNewStat()
	return nil
}

func (s *StatsView) updateTimer(t *Timer, g *gocui.Gui) {
	defer t.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.getDoneCh():
			return
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				// TODO remove hardcoded view name
				v, err := g.View("stats")
				if err != nil {
					return err
				}
				v.Clear()
				elapsed, _ := t.ElapsedTime()
				fmt.Fprintf(v, "%02d:%02d", elapsed.Mins, elapsed.Secs)
				return nil
			})
		}
	}
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

// SetKeyBindings adds keybindings to scroll in stats history
func (s *StatsView) SetKeyBindings(g *gocui.Gui) {
	if err := g.SetKeybinding(s.name, gocui.KeyCtrlJ, gocui.ModNone, s.stats.ScrollUp); err != nil {
		Logger.Warn(fmt.Sprintf("%v", err))
	}
	if err := g.SetKeybinding(s.name, gocui.KeyCtrlK, gocui.ModNone, s.stats.ScrollDown); err != nil {
		Logger.Warn(fmt.Sprintf("%v", err))
	}
}

// UnsetKeyBindings deletes keybindigns when this view is active
func (s *StatsView) UnsetKeyBindings(g *gocui.Gui) {
	g.DeleteKeybindings(s.name)
}
