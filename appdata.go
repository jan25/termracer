package main

import (
	"github.com/jan25/gocui"
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

// AppData wraps all view's data structs in the app
type AppData struct {
	paragraph *viewdata.ParagraphData
	editor    *viewdata.WordEditorData
	history   *viewdata.Stats
	stats     *viewdata.LiveStats
	controls  *viewdata.ControlsData
}

// InitializeAppData creates views and initialises AppData
func InitializeAppData(g *gocui.Gui) (*AppData, error) {
	ad := &AppData{}

	para := views.NewParagraphView(paraName, topX, topY, paraW, paraH)
	ad.paragraph = para.Data

	word := views.NewWordView(wordName, topX, topY+paraH+pad, wordW, wordH)
	ad.editor = word.Data

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
	updateUICh := make(chan bool)
	paraToWord := make(chan viewdata.WordValidateMsg)
	wordToPara := make(chan viewdata.WordValidateMsg)
	wordToStats := make(chan viewdata.StatMsg)
	ad.paragraph.SetChannels(paraToWord, wordToPara, updateUICh)
	ad.editor.SetChannels(wordToPara, paraToWord, wordToStats, updateUICh)
	ad.stats.SetChannels(wordToStats, updateUICh)

	ad.stats.IsActive = !ad.stats.IsActive
	ad.history.IsActive = !ad.history.IsActive

	if err := ad.paragraph.StartRace(); err != nil {
		return err
	}
	if err := ad.editor.StartRace(g, wordName); err != nil {
		return err
	}
	if err := ad.stats.StartRace(); err != nil {
		return err
	}
	ad.controls.StartRace()

	return nil
}

// OnRaceFinish at end of a race:
// - when typing is finished
// - when user stops the race
func (ad *AppData) OnRaceFinish(g *gocui.Gui) error {
	if err := ad.paragraph.FinishRace(); err != nil {
		return err
	}
	if err := ad.editor.FinishRace(g); err != nil {
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
