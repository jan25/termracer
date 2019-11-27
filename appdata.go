package main

import (
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/viewdata"
	"github.com/jan25/termracer/views"
)

// AppData wraps all view's data structs in the app
type AppData struct {
	paragraph *viewdata.ParagraphData
	editor    *viewdata.WordEditorData
	history   *viewdata.Stats
	stats     *viewdata.LiveStats
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

	g.SetManager(para, word, stats, controls)

	return ad, nil
}

// OnRaceStart is called at start of a new race
func (ad *AppData) OnRaceStart() error {
	paraToWord := make(chan viewdata.WordValidateMsg)
	wordToPara := make(chan viewdata.WordValidateMsg)
	wordToStats := make(chan viewdata.StatMsg)
	ad.paragraph.SetChannels(paraToWord, wordToPara)
	ad.editor.SetChannels(wordToPara, paraToWord, wordToStats)
	ad.stats.SetChannels(wordToStats)

	if err := ad.paragraph.StartRace(); err != nil {
		return err
	}
	if err := ad.editor.StartRace(); err != nil {
		return err
	}
	if err := ad.stats.StartRace(); err != nil {
		return err
	}
	// TODO: do the above for controls view

	ad.stats.IsActive = !ad.stats.IsActive
	ad.history.IsActive = !ad.history.IsActive

	return nil
}

// OnRaceFinish at end of a race:
// - when typing is finished
// - when user stops the race
func (ad *AppData) OnRaceFinish() error {
	if err := ad.paragraph.FinishRace(); err != nil {
		return err
	}
	if err := ad.editor.FinishRace(); err != nil {
		return err
	}
	if err := ad.stats.FinishRace(); err != nil {
		return err
	}
	// TODO: finish race for controls too

	ad.stats.IsActive = !ad.stats.IsActive
	ad.history.IsActive = !ad.history.IsActive

	return nil
}
