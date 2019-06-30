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

	// content to display
	content string
}

func newControls(name string, x, y int, w, h int) *Controls {
	return &Controls{
		name:    name,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		content: defaultContent,
	}
}

// Layout manager for controls View
func (c *Controls) Layout(g *gocui.Gui) error {
	v, err := g.SetView(c.name, c.x, c.y, c.x+c.w, c.y+c.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Title = title

	v.Clear()
	fmt.Fprintf(v, "%s", c.content)
	return nil
}

// RaceModeControls sets up controls
// for during a race
func (c *Controls) RaceModeControls() {
	c.content = raceModeContent
}

// DefaultControls sets up controls
// when no race in progress
func (c *Controls) DefaultControls() {
	c.content = defaultContent
}

const title = "Controls"

const defaultContent = `
Ctrl+s		New race
Ctrl+j/k	Scroll
Ctrl+c 		Quit`

const raceModeContent = `
Ctrl+e End race
Ctrl+c Quit`
