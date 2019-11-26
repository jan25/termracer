package views

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/pkg/utils"
	"github.com/jan25/termracer/viewdata"
)

// StatsView is to keep track of
// stats and the view
type StatsView struct {
	// name of the view
	name string

	// position and dims
	x, y int
	w, h int

	// Stats history data
	HistoryData *viewdata.Stats

	// Live Stats for race in progress
	LiveRaceData *viewdata.LiveStats
}

// NewStatsView creates brand new StatsView instance
func NewStatsView(name string, x, y, w, h int) (*StatsView, error) {
	hd, err := viewdata.NewStatsHistory()
	if err != nil {
		return nil, err
	}
	return &StatsView{
		name:         name,
		x:            x,
		y:            y,
		w:            w,
		h:            h,
		HistoryData:  hd,
		LiveRaceData: viewdata.NewLiveStats(),
	}, nil
}

// Layout manager for Stats View
func (sv *StatsView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(sv.name, sv.x, sv.y, sv.x+sv.w, sv.y+sv.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if sv.HistoryData.IsActive {
		return sv.renderHistoryData(v)
	} else if sv.LiveRaceData.IsActive {
		return sv.renderLiveRaceData(v, g)
	} else {
		return errors.New("Failed to render stats view content")
	}
}

func (sv *StatsView) renderHistoryData(v *gocui.View) error {
	v.Title = "Race History"
	v.Clear()

	// FIXME: Any error case to handle?

	history := sv.HistoryData

	if len(history.List) == 0 {
		v.Wrap = true
		cyan := color.New(color.FgCyan)
		cyan.Fprintf(v, "No races available to show")
		return nil
	}

	// Heading
	ul := color.New(color.Underline)
	ul.Fprintln(v, "wpm acc. when    ")

	selected := history.Selected
	for i := selected; i >= 0; i-- {
		stat := history.List[i]
		f := "%s\n"
		if i == history.Selected {
			white := color.New(color.BgWhite)
			white.Fprintf(v, f,
				fmt.Sprintf("%-3d %3d%% %-8s", stat.Wpm, int(stat.Accuracy), utils.FormatDate(stat.When)))
		} else {
			fmt.Fprintf(v, f,
				fmt.Sprintf("%-3d %3d%% %-8s", stat.Wpm, int(stat.Accuracy), utils.FormatDate(stat.When)))
		}
	}

	return nil
}

func (sv *StatsView) renderLiveRaceData(v *gocui.View, g *gocui.Gui) error {
	sv.LiveRaceData.TryStartTicker(g)

	v.Title = "Race in progress"
	v.Clear()

	t, err := sv.LiveRaceData.ElapsedTime()
	if err != nil {
		return err
	}
	fmt.Fprintf(v, "%02d:%02d \n", t.Mins, t.Secs)

	wpm, err := sv.LiveRaceData.Wpm()
	if err != nil {
		return err
	}
	acc, err := sv.LiveRaceData.Accuracy()
	if err != nil {
		return err
	}
	fmt.Fprintf(v, "wpm %d \nAccuracy %.2f%% \n", wpm, acc)

	return nil
}
