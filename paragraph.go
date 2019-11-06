package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jan25/color"
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
	// position values of view
	pos ViewPosition
	// whether current word is mistyped
	Mistyped bool
	// count of lines for a paragraph sample
	lineCount int
}

// ViewPosition keeps track of top left Origin of View
// and also line, word number of highlighted word
type ViewPosition struct {
	// Y position of View origin
	Oy int

	// For highlighted word: Line, Word number both start at 0
	Line int
	Word int
}

func newParagraph(name string, x, y int, w, h int) *Paragraph {
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
	for i, w := range p.words {
		p.words[i] = strings.TrimSpace(w)
	}

	n := AddNewLines(p.words, p.w-1)
	p.lineCount = n
	p.wordi = 0
	p.pos = ViewPosition{
		Oy:   0,
		Line: 0,
		Word: 0,
	}
}

// Advance moves target word to next word
func (p *Paragraph) Advance() error {
	if p.wordi >= len(p.words)-1 {
		return errors.New("can not advance beyond number of words")
	}

	p.wordi++

	// Scroll logic
	prevWord := p.words[p.wordi-1]
	if strings.HasSuffix(prevWord, "\n") {
		// Jump to next line
		p.pos.Line++
		p.pos.Word = 0
	} else {
		p.pos.Word++
	}
	p.makeScroll()

	return nil
}

// CountDoneWords returns number of words
// already types in current race
func (p *Paragraph) CountDoneWords() int {
	return p.wordi
}

// CurrentWord returns target word to type
func (p *Paragraph) CurrentWord() string {
	w := p.words[p.wordi]
	if strings.HasSuffix(w, "\n") {
		return strings.TrimRight(w, "\n")
	}
	return w
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

	v.SetOrigin(0, p.pos.Oy) // scroll

	for i, w := range p.words {
		highlight := false
		done := false
		if i < p.wordi {
			done = true
		} else if i == p.wordi {
			highlight = true
		}

		p.printWord(v, w, highlight, done)
	}
}

func (p *Paragraph) printWord(v *gocui.View, w string, highlight bool, done bool) {
	f := "%s"
	if done {
		fg := color.New(color.FgGreen)
		fg.Fprintf(v, f, w)
	} else if highlight {
		bg := color.New(color.BgGreen)
		if p.Mistyped {
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

func (p *Paragraph) makeScroll() {
	atLinesLeft := 2
	atWord, atLine := 2, (p.h-1)-atLinesLeft

	if p.pos.Word != atWord {
		return
	}
	currLine := p.pos.Line - p.pos.Oy
	linesLeft := (p.lineCount - 1) - p.pos.Line
	if currLine == atLine && linesLeft >= atLinesLeft {
		p.pos.Oy++
	}
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

	p.Mistyped = false
}
