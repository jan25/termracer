package main

import (
	"log"

	"github.com/jan25/gocui"
)

const (
	STATS_VIEW    = "stats"
	PARA_VIEW     = "para"
	WORD_VIEW     = "word"
	CONTROLS_VIEW = "controls"
)

var (
	g         *gocui.Gui
	paragraph *Paragraph
	word      *Word
	stats     *Stats
)

var (
	paraW, paraH = 60, 8
	wordW, wordH = 60, 2

	statsW, statsH       = 20, 6
	controlsW, controlsH = 20, 4

	topX, topY = 1, 1
	pad        = 1
)

func main() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g = gui
	defer g.Close()

	paragraph = newParagraph(PARA_VIEW, topX, topY, paraW, paraH)
	word = newWord(WORD_VIEW, topX, topY+paraH+pad, wordW, wordH)
	stats = newStatsView(STATS_VIEW, topX+paraW+pad, topY, statsW, statsH)
	controls := newControls(CONTROLS_VIEW, topX+paraW+pad, topY+statsH+pad, controlsW, controlsH)

	g.SetManager(paragraph, word, stats, controls)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, ctrlS); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone, ctrlE); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func ctrlS(g *gocui.Gui, v *gocui.View) error {
	paragraph.Init()
	word.Init()
	stats.StartTimer()

	g.SetCurrentView(WORD_VIEW)
	g.Cursor = true

	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	paragraph.Reset()
	word.Reset()
	stats.StopTimer()

	g.Cursor = false

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
