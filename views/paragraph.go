package views

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jan25/gocui"
	viewdata "github.com/jan25/termracer/views/data"
)

// ParagraphView keeps track of
// its view's metadata
type ParagraphView struct {
	// name of View
	name string
	// positions, dimensions
	x, y int
	w, h int

	// done channel
	done chan struct{}

	// Data stores content of the view
	Data *viewdata.ParagraphData
}

// NewParagraphView creates new instance of ParagraphView
func NewParagraphView(name string, x, y int, w, h int) *ParagraphView {
	pd := viewdata.NewParagraphData()
	pd.H = h
	pd.W = w
	return &ParagraphView{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
		Data: pd,
	}
}

// Layout manager for ParagraphView
func (pv *ParagraphView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(pv.name, pv.x, pv.y, pv.x+pv.w, pv.y+pv.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	select {
	case <-pv.Data.DoneCh():
		// no race in progress
		v.Clear()
	default:
		pv.DrawView(v)
	}

	return nil
}

// DrawView renders the ParagraphView
func (pv *ParagraphView) DrawView(v *gocui.View) {
	v.Clear()

	v.SetOrigin(0, pv.Data.Oy) // scroll

	pd := pv.Data
	for i, w := range pd.Words {
		highlight := false
		done := false

		ti := pd.GetCurrentIdx()
		if i < ti {
			done = true
		} else if i == ti {
			highlight = true
		}

		pv.printWord(v, w, highlight, done)
	}
}

func (pv *ParagraphView) printWord(v *gocui.View, w string, highlight bool, done bool) {
	f := "%s"
	if done {
		fg := color.New(color.FgGreen)
		fg.Fprintf(v, f, w)
	} else if highlight {
		bg := color.New(color.BgGreen)
		if pv.Data.Mistyped {
			bg = color.New(color.BgRed)
		}
		bg.Add(color.Underline)
		bg.Fprintf(v, f, w)
	} else {
		fmt.Fprintf(v, f, w)
	}

	if !strings.HasSuffix(w, "\n") {
		// Space between words in a paragraph
		fmt.Fprint(v, " ")
	}
}
