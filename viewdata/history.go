package viewdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/pkg/utils"
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
	List     []*OneStat
	Selected int // To keep track of highlighted stat in history

	IsActive bool // whether this data is active to show in stats view
}

// NewStatsHistory creates new Stats instance to keep track of race history
func NewStatsHistory() (*Stats, error) {
	s := Stats{
		IsActive: true, // default data on app startup
	}
	if err := s.LoadHistory(); err != nil {
		return nil, err
	}
	return &s, nil
}

// Just use s.IsActive directly for now
// instead of:
// OnStartRace called at start of new race
// func (s *Stats) OnStartRace() {
// 	s.IsActive = false
// }
// OnFinishRace called at race finish
// func (s *Stats) OnFinishRace() {
// 	s.IsActive = true
// }

// SaveNewStat saves current race's stat to history of stats
func (s *Stats) SaveNewStat(stat *OneStat) error {
	s.List = append(s.List, stat)
	s.appendToFile(stat) // save to history file store
	// Always selected most recent race stat on finish
	s.Selected = len(s.List) - 1
	return nil
}

func (s *Stats) appendToFile(stat *OneStat) error {
	line := fmt.Sprintf("%d,%f,%s", stat.Wpm, stat.Accuracy, utils.FormatDate(stat.When))
	f, _ := config.GetHistoryFilePath()
	if err := utils.AppendLineEOF(f, line); err != nil {
		// FIXME add logs if needed as below
		// Logger.Warn("Failed to append to file", zap.Error(err))
		return err
	}
	// FIXME add logs if needed as below
	// Logger.Info("Successfuly appended to file")
	return nil
}

// LoadHistory loads history data from local FS
func (s *Stats) LoadHistory() error {
	fname, _ := config.GetHistoryFilePath()
	f, err := os.Open(fname)
	if err != nil {
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
			if s.List == nil {
				s.List = make([]*OneStat, 0)
			}
			wpm, err := strconv.Atoi(record[cmapper["wpm"]])
			if err != nil {
				// FIXME add a log for below case
				// Logger.Warn("failed to read a record: " + strings.Join(record, ""))
				continue
			}
			acc, err := strconv.ParseFloat(record[cmapper["acc"]], 64)
			if err != nil {
				// FIXME add a log for below case
				// Logger.Warn("failed to read a record: " + strings.Join(record, ""))
				continue
			}
			when, err := time.Parse("02/01/06", record[cmapper["when"]])

			stat := &OneStat{
				Wpm:      wpm,
				Accuracy: acc,
				When:     when,
			}
			s.List = append(s.List, stat)
		}
	}
	// FIXME add logger info if needed as below
	// Logger.Info("Finished loading records from file", zap.Int("records", len(records)))
	s.Selected = len(s.List) - 1

	return nil
}

// ScrollDown is a keybinding
// increments selected stat index
func (s *Stats) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	if s.Selected+1 < len(s.List) {
		s.Selected++
	} else {
		config.Logger.Info("End of history reached. Can not scroll further up")
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
		config.Logger.Info("Top of history reached. Can not scroll further down")
	}
	// just to adhere to KeyBinding handler interface
	return nil
}
