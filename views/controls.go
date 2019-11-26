package views

import (
	"fmt"

	"github.com/jan25/gocui"
)

// ControlsView represents static
// view to display controls
type ControlsView struct {
	// name of View
	name string

	// position and dims
	x, y int
	w, h int

	// content to display
	content string
}

// NewControls creates new instance of ControlsView
func NewControls(name string, x, y int, w, h int) *ControlsView {
	return &ControlsView{
		name:    name,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		content: defaultContent,
	}
}

// Layout manager for controls View
func (c *ControlsView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(c.name, c.x, c.y, c.x+c.w, c.y+c.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Title = title

	v.Clear()
	fmt.Fprintf(v, "%s", c.content)
	return nil
}

// RaceModeControlsContent sets up controls
// for during a race
func (c *ControlsView) RaceModeControlsContent() {
	c.content = raceModeContent
}

// DefaultControlsContent sets up controls
// when no race in progress
func (c *ControlsView) DefaultControlsContent() {
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
