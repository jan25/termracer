package main

import (
	"log"

	"github.com/jan25/gocui"
	"go.uber.org/zap"
)

const (
	statsName    = "stats"
	paraName     = "para"
	wordName     = "word"
	controlsName = "controls"
)

var (
	logger    *zap.Logger
	g         *gocui.Gui
	paragraph *Paragraph
	word      *Word
	stats     *StatsView
)

var (
	paraW, paraH = 60, 8
	wordW, wordH = 60, 2

	statsW, statsH       = 20, 6
	controlsW, controlsH = 20, 4

	topX, topY = 1, 1
	pad        = 1
)

func initLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./logs/app.log",
	}
	logger, _ = cfg.Build()
}

func main() {
	initLogger()
	defer logger.Sync()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g = gui
	defer g.Close()

	paragraph = newParagraph(paraName, topX, topY, paraW, paraH)
	word = newWord(wordName, topX, topY+paraH+pad, wordW, wordH)
	stats = newStatsView(statsName, topX+paraW+pad, topY, statsW, statsH)
	controls := newControls(controlsName, topX+paraW+pad, topY+statsH+pad, controlsW, controlsH)

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
	stats.StartRace()

	g.SetCurrentView(wordName)
	g.Cursor = true

	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	paragraph.Reset()
	word.Reset()
	stats.StopRace()

	g.Cursor = false

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
