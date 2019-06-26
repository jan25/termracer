package main

import (
	"strings"

	"github.com/jan25/gocui"
)

const MAX_WORD_LEN int = 15

// for developerment
// remove when this module is fully developed
// const TEST_WORD string = "testword"

type Word struct {
	// name of the View
	name string
	// position, dimentions
	x, y int
	w, h int

	// editor instance
	e gocui.Editor
	// keep track of status of view
	done chan struct{}
}

func newWord(name string, x, y int, w, h int) *Word {
	return &Word{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
	}
}

func (w *Word) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	select {
	case <-w.done:
		// channel closed
		v.Clear()
	default:
		w.init(v)
	}
	return nil
}

func (w *Word) Init() {
	w.done = make(chan struct{})
}

func (w *Word) Reset() {
	close(w.done)
}

func (w *Word) init(v *gocui.View) {
	w.e = newWordEditor()

	v.Editor = newWordEditor()
	v.Editable = true
	v.SelBgColor = gocui.ColorRed
	v.SelFgColor = gocui.ColorCyan
}

func newWordEditor() gocui.Editor {
	return gocui.EditorFunc(wordEditorFunc)
}

func wordEditorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	word := strings.TrimSpace(getCurrentWord(v))

	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		handleDelete(v, true)
	case key == gocui.KeyDelete:
		handleDelete(v, false)
	case len(word) > MAX_WORD_LEN:
		// do not add anymore runes
		// can only delete from here
	case ch != 0 && mod == 0:
		handleChar(v, ch)
	case key == gocui.KeySpace:
		handleSpace(v)
	}
}

func handleDelete(v *gocui.View, back bool) {
	v.EditDelete(back)
	checkAndHighlight(v)
}

func handleChar(v *gocui.View, ch rune) {
	v.EditWrite(ch)
	checkAndHighlight(v)
}

func checkAndHighlight(v *gocui.View) {
	w := strings.TrimSpace(getCurrentWord(v))
	ok := strings.HasPrefix(paragraph.CurrentWord(), w)
	highlight(ok, v)
}

func handleSpace(v *gocui.View) {
	w := strings.TrimSpace(getCurrentWord(v))
	if w == paragraph.CurrentWord() {
		v.Clear()
		v.SetCursor(v.Origin())

		perr := paragraph.Advance()
		if perr != nil {
			paragraph.Reset()
			word.Reset()
		}
		// paragraph.DrawView()
	} else {
		highlight(false, v)
		v.EditWrite(' ')
	}
}

func highlight(ok bool, v *gocui.View) {
	v.Highlight = !ok
}

func getCurrentWord(v *gocui.View) string {
	line := v.Buffer()
	return line
}
