package main

import (
	"fmt"

	"github.com/jan25/gocui"
)

type Controls struct {
	// name of View
	name string

	// position and dims
	x, y int
	w, h int
}

func newControls(name string, x, y int, w, h int) *Controls {
	return &Controls{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
	}
}

func (c *Controls) Layout(g *gocui.Gui) error {
	v, err := g.SetView(c.name, c.x, c.y, c.x+c.w, c.y+c.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	fmt.Fprintf(v, "%s", CONTROLS)
	return nil
}

const CONTROLS = `
Ctrl+c Quit
Ctrl+s New race
Ctrl+e End race`
