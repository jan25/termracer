package views

import (
	"fmt"

	"github.com/jan25/gocui"
	viewdata "github.com/jan25/termracer/views/data"
)

// ControlsView represents static
// view to display controls
type ControlsView struct {
	// name of View
	name string

	// position and dims
	x, y int
	w, h int

	Data *viewdata.ControlsData
}

// NewControls creates new instance of ControlsView
func NewControls(name string, x, y int, w, h int) *ControlsView {
	return &ControlsView{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
		Data: viewdata.NewControlsData(),
	}
}

// Layout manager for controls View
func (c *ControlsView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(c.name, c.x, c.y, c.x+c.w, c.y+c.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Title = "Controls"

	v.Clear()
	fmt.Fprintf(v, "%s", c.Data.Content)
	return nil
}
