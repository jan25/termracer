package views

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/viewdata"
)

// ParagraphView keeps track of
// its view's metadata
type ParagraphView struct {
	// name of View
	name string
	// positions, dimensions
	x, y int
	w, h int

	// Y position of View origin
	Oy int

	// done channel
	done chan struct{}

	// Data stores content of the view
	Data *viewdata.ParagraphData
}

func newParagraphView(name string, x, y int, w, h int) *ParagraphView {
	return &ParagraphView{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
		Oy:   0,
		Data: viewdata.NewParagraphData(),
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

	v.SetOrigin(0, pv.Oy) // scroll

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

// FIXME unused method. Maybe call it in Layout() ?
func (pv *ParagraphView) makeScroll() {
	whenLinesLeft := 2
	atWord, atLine := 2, (pv.h-1)-whenLinesLeft

	pd := pv.Data
	if pd.Word != atWord {
		return
	}
	currLine := pd.Line - pv.Oy
	linesLeft := (pd.GetLineCount() - 1) - pd.Line
	if currLine == atLine && linesLeft >= whenLinesLeft {
		pv.Oy++
	}
}
