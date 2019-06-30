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
	// name of View
	name string
	// positions, dimensions
	x, y int
	w, h int

	// done channel
	done chan struct{}
	// words in pargraph
	words []string
	// index of current word being typed
	wordi int
}

func newParagraph(name string, x, y int, w, h int) *Paragraph {
	// split into words at whitespace characters
	// words := strings.Fields(paragraph)

	// view.Wrap = true
	return &Paragraph{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
	}
}

// Layout manager for paragraph View
func (p *Paragraph) Layout(g *gocui.Gui) error {
	v, err := g.SetView(p.name, p.x, p.y, p.x+p.w, p.y+p.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Wrap = true

	select {
	case <-p.getDoneCh():
		// channel closed
		v.Clear()
	default:
		p.DrawView(v)
	}

	return nil
}

// Init initialises new paragraph to type
func (p *Paragraph) Init() {
	p.done = make(chan struct{})

	para, err := ChooseParagraph()
	if err != nil {
		Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	p.words = strings.Fields(para)
	p.wordi = 0
}

// Advance moves target word to next word
func (p *Paragraph) Advance() error {
	if p.wordi >= len(p.words)-1 {
		return errors.New("can not advance beyond number of words")
	}

	p.wordi++
	return nil
}

// CountDoneWords returns number of words
// already types in current race
func (p *Paragraph) CountDoneWords() int {
	return p.wordi
}

// CurrentWord returns target word to type
func (p *Paragraph) CurrentWord() string {
	return p.words[p.wordi]
}

// CharsUptoCurrent counts chars in words
// upto/including current word
func (p *Paragraph) CharsUptoCurrent() int {
	c := 0
	for i := 0; i < p.wordi; i++ {
		c += len(p.words[i])
	}
	return c
}

// DrawView renders the paragraph View
func (p *Paragraph) DrawView(v *gocui.View) {
	v.Clear()

	for i, w := range p.words {
		underline := false
		if i == p.wordi {
			underline = true
		}

		p.printWord(v, w, underline)
	}
}

func (p *Paragraph) printWord(v *gocui.View, w string, underline bool) {
	f := "%s "
	if underline {
		f = "\033[0;7m%s\033[0m "
	}
	fmt.Fprintf(v, f, w)
}

func (p *Paragraph) getDoneCh() chan struct{} {
	if p.done == nil {
		p.done = make(chan struct{})
	}
	return p.done
}

// Reset deactivates the paragraph view
// used to stop a race
func (p *Paragraph) Reset() {
	select {
	case <-p.getDoneCh():
		// already closed
		// nothing to do
	default:
		close(p.getDoneCh())
	}
}
