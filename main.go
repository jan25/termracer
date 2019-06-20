package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

var timer *time.Timer = time.NewTimer(time.Millisecond)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
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
		b, err := ioutil.ReadFile("sample_paragraph.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(para, "%s", b)

		para.Wrap = true

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if word, err := g.SetView("word",
		topX, topY+paraH+pad,
		topX+wordW, topY+paraH+pad+wordH); err != nil {

		word.Editable = true

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	g.SetCurrentView("word")
	g.Cursor = true

	if stats, err := g.SetView("stats",
		topX+paraW+pad, topY,
		topX+paraW+pad+statsW, topY+statsH); err != nil {

		time := <-timer.C
		fmt.Fprintf(stats, "%v", time)

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
