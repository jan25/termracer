package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/jan25/gocui"
	// "go.uber.org/zap"
)

var (
	// logger    zap.Logger
	paragraph *Paragraph
)

var (
	done = make(chan struct{})
	wg   sync.WaitGroup
)

func main() {
	// logger, _ := zap.NewProduction()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	wg.Add(1)
	go updateTimer(g)

	// logger.Info("started gui..")
	// defer logger.Sync()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
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

	if stats, err := g.SetView("stats",
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

func updateTimer(g *gocui.Gui) {
	defer wg.Done()

	timer := NewTimer()
	timer.Start()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("stats")
				if err != nil {
					return err
				}
				v.Clear()
				elapsed, _ := timer.ElapsedTime()
				fmt.Fprintf(v, "%02d:%02d", elapsed.Mins, elapsed.Secs)
				return nil
			})
		}
	}
}
