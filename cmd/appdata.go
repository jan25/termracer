package main

import (
	"log"
	
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/viewdata"
	"github.com/jan25/termracer/views"
)

const (
	statsName    = "stats"
	paraName     = "para"
	wordName     = "word"
	controlsName = "controls"
)

var (
	paraW, paraH = 60, 8
	wordW, wordH = 60, 2

	statsW, statsH       = 20, 6
	controlsW, controlsH = 20, 4

	topX, topY = 1, 1
	pad        = 1
)

var app *AppData

// AppData wraps all view's data structs in the app
type AppData struct {
	paragraph *viewdata.ParagraphData
	history   *viewdata.Stats
	stats     *viewdata.LiveStats
	controls  *viewdata.ControlsData

	updateUICh chan bool
	finishCh   chan viewdata.OneStat
}

// initializeAppData creates views and initialises AppData
func initializeAppData(g *gocui.Gui) (*AppData, error) {
	ad := &AppData{}

	para := views.NewParagraphView(paraName, topX, topY, paraW, paraH)
	ad.paragraph = para.Data

	word := views.NewWordView(wordName, topX, topY+paraH+pad, wordW, wordH)
	word.Data = para.Data // This Data shared between editor and paragraph views

	stats, err := views.NewStatsView(statsName, topX+paraW+pad, topY, statsW, statsH)
	if err != nil {
		return nil, err
	}
	ad.stats = stats.LiveRaceData
	ad.history = stats.HistoryData

	controls := views.NewControls(controlsName, topX+paraW+pad, topY+statsH+pad, controlsW, controlsH)
	ad.controls = controls.Data

	g.SetManager(para, word, stats, controls)

	return ad, nil
}

// OnRaceStart is called at start of a new race
func (ad *AppData) OnRaceStart(g *gocui.Gui) error {
	ad.updateUICh = make(chan bool)            // close()ed in ticker.go
	paraToStats := make(chan viewdata.StatMsg) // close()ed in livestats.go
	ad.finishCh = make(chan viewdata.OneStat)
	ad.stats.SetChannels(paraToStats, ad.updateUICh, ad.finishCh)
	ad.paragraph.SetChannels(paraToStats, ad.updateUICh)

	ad.stats.IsActive = !ad.stats.IsActive
	ad.history.IsActive = !ad.history.IsActive

	ad.paragraph.StartRace(g, wordName)
	if err := ad.stats.StartRace(); err != nil {
		return err
	}
	ad.controls.StartRace()

	// Wait for end of race signal
	// from ParagraphData
	go func(g *gocui.Gui) {
		for newStat := range ad.finishCh {
			ad.history.SaveNewStat(&newStat)
			// Below two are similar to keybindings.ctrlE func
			err := ad.OnRaceFinish()
			afterRaceControls(g, false)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}(g)

	return nil
}

// OnRaceFinish at end of a race:
// - when typing is finished
// - when user stops the race
func (ad *AppData) OnRaceFinish() error {
	config.Logger.Info("OnRaceFinish()")
	close(ad.finishCh)

	if err := ad.paragraph.FinishRace(); err != nil {
		return err
	}
	if err := ad.stats.FinishRace(); err != nil {
		return err
	}
	ad.controls.DefaultControls()

	ad.stats.IsActive = !ad.stats.IsActive
	ad.history.IsActive = !ad.history.IsActive

	return nil
}

// HistoryScrollUp scrolls the start history list
func (ad *AppData) HistoryScrollUp(g *gocui.Gui, v *gocui.View) error {
	return ad.history.ScrollUp(g, v)
}

// HistoryScrollDown scrolls the start history list
func (ad *AppData) HistoryScrollDown(g *gocui.Gui, v *gocui.View) error {
	return ad.history.ScrollDown(g, v)
}

// DebugAdvance is used to debug manually
func (ad *AppData) DebugAdvance(g *gocui.Gui, v *gocui.View) error {
	ad.paragraph.DebugAdvance()
	return nil
}
