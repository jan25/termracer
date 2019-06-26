package main

import (
	"fmt"

	"github.com/jan25/gocui"
)

// Controls represents static
// view to display controls
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

// Layout manager for controls View
func (c *Controls) Layout(g *gocui.Gui) error {
	v, err := g.SetView(c.name, c.x, c.y, c.x+c.w, c.y+c.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	fmt.Fprintf(v, "%s", controlsContent)
	return nil
}

const controlsContent = `
Ctrl+c Quit
Ctrl+s New race
Ctrl+e End race`
