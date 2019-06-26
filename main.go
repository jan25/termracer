package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jan25/gocui"
)

const (
	STATS_VIEW = "stats"
	PARA_VIEW  = "para"
	WORD_VIEW  = "word"
)

var (
	g         *gocui.Gui
	paragraph *Paragraph
	word      *Word
	timer     *Timer
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

	timer = NewTimer()

	paragraph = newParagraph(PARA_VIEW, topX, topY, paraW, paraH)
	word = newWord(WORD_VIEW, topX, topY+paraH+pad, wordW, wordH)

	g.SetManager(paragraph, word)
	// g.SetManagerFunc(layout)

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
	// timer.Start()
	paragraph.Init()
	word.Init()

	g.SetCurrentView(WORD_VIEW)
	g.Cursor = true

	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	// timer.Stop()
	paragraph.Reset()
	word.Reset()
	g.Cursor = false

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	// timer.Stop()

	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	// maxX, maxY := g.Size()

	if stats, err := g.SetView(STATS_VIEW,
		topX+paraW+pad, topY,
		topX+paraW+pad+statsW, topY+statsH); err != nil {

		fmt.Fprintf(stats, "%02d:%02d", 10, 1)

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if controls, err := g.SetView("controls",
		topX+paraW+pad, topY+statsH+pad,
		topX+paraW+pad+controlsW, topY+statsH+pad+controlsH); err != nil {

		controls.Title = "Controls"

		b, err := ioutil.ReadFile("controls.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(controls, "%s", b)

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}
