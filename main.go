package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jan25/gocui"
	// "go.uber.org/zap"
)

const (
	STATS_VIEW = "stats"
)

var (
	// logger    zap.Logger
	g         *gocui.Gui
	paragraph *Paragraph
	timer     *Timer
)

func main() {
	// logger, _ := zap.NewProduction()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g = gui
	defer g.Close()

	timer = NewTimer()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, ctrlS); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone, ctrlE); err != nil {
		log.Panicln(err)
	}

	// logger.Info("started gui..")
	// defer logger.Sync()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func ctrlS(g *gocui.Gui, v *gocui.View) error {
	timer.Start()

	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	timer.Stop()

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	timer.Stop()

	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	// maxX, maxY := g.Size()

	paraW, paraH := 60, 8
	wordW, wordH := 60, 2

	statsW, statsH := 20, 6
	controlsW, controlsH := 20, 4

	topX, topY := 1, 1
	pad := 1

	if para, err := g.SetView("para", topX, topY, topX+paraW, topY+paraH); err != nil {
		b, err := ioutil.ReadFile("samples/sample_paragraph.txt")
		if err != nil {
			panic(err)
		}

		paragraph = newParagraph(string(b), para)
		paragraph.DrawView()

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if word, err := g.SetView("word",
		topX, topY+paraH+pad,
		topX+wordW, topY+paraH+pad+wordH); err != nil {

		word.Editor = WordEditor
		word.Editable = true
		word.SelBgColor = gocui.ColorRed
		word.SelFgColor = gocui.ColorCyan

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	g.SetCurrentView("word")
	g.Cursor = true

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
