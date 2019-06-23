package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jan25/gocui"
)

// Paragraph encapsulates the data and
// state of Paragraph view
type Paragraph struct {
	// words in pargraph
	words []string

	// index of current word being typed
	currentw int

	// instance of view to display paragraph
	view *gocui.View
}

func newParagraph(paragraph string, view *gocui.View) *Paragraph {
	// split into words at whitespace characters
	words := strings.Fields(paragraph)

	view.Wrap = true
	return &Paragraph{
		words:    words,
		currentw: 0,
		view:     view,
	}
}

func (p *Paragraph) Advance() error {
	if p.currentw >= len(p.words) {
		return errors.New("can not advance beyond number of words")
	}

	p.currentw += 1
	return nil
}

func (p *Paragraph) DrawView() {
	p.view.Clear()

	for i, w := range p.words {
		underline := false
		if i == p.currentw {
			underline = true
		}

		p.printWord(w, underline)
	}
}

func (p *Paragraph) printWord(w string, underline bool) {
	f := "%s "
	if underline {
		f = "\033[0;4m%s\033[0m "
	}
	fmt.Fprintf(p.view, f, w)
}

// for development purpose
// remove once this module is fairly integrated
func (p *Paragraph) Reset() {
	p.currentw = 0
}
