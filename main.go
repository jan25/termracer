
package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

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
	wordW, wordH := 60, 3

	statsW, statsH := 20, 6
	controlsW, controlsH := 20, 5

	topX, topY := 1, 1
	pad := 1

	if _, err := g.SetView("para", topX, topY, topX + paraW, topY + paraH); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView("word",
		topX, topY + paraH + pad,
		topX + wordW, topY + paraH + pad + wordH); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView("stats",
		topX + paraW + pad, topY,
		topX + paraW + pad + statsW, topY + statsH); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView("controls",
		topX + paraW + pad, topY + statsH + pad,
		topX + paraW + pad + controlsW, topY + statsH + pad + controlsH); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
	}
	
	return nil
}