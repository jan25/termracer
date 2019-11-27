package main

// import (
// 	"strings"

// 	"github.com/jan25/gocui"
// )

// const maxWordLen int = 15

// // Word is a widget to type
// // a target word
// type Word struct {
// 	// name of the View
// 	name string
// 	// position, dimentions
// 	x, y int
// 	w, h int

// 	// editor instance
// 	e gocui.Editor
// 	// keep track of status of view
// 	done chan struct{}

// 	// mistakes done during race
// 	Mistyped int
// }

// func newWord(name string, x, y int, w, h int) *Word {
// 	return &Word{
// 		name: name,
// 		x:    x,
// 		y:    y,
// 		w:    w,
// 		h:    h,
// 	}
// }

// // Layout is layout manager for word widget
// func (w *Word) Layout(g *gocui.Gui) error {
// 	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
// 	if err != nil && err != gocui.ErrUnknownView {
// 		return err
// 	}

// 	select {
// 	case <-w.getDoneCh():
// 		// channel closed
// 		clearEditor(v)
// 	default:
// 		w.init(v)
// 	}
// 	return nil
// }

// // Init initialises the word View
// // This makes the widget editable
// func (w *Word) Init() {
// 	w.done = make(chan struct{})

// 	g.SetCurrentView(wordName)
// 	g.Cursor = true
// }

// func (w *Word) getDoneCh() chan struct{} {
// 	if w.done == nil {
// 		w.done = make(chan struct{})
// 	}
// 	return w.done
// }

// // Reset clears the widget
// // used when race is not active
// func (w *Word) Reset() {
// 	select {
// 	case <-w.getDoneCh():
// 		// already closed
// 		// nothing to do
// 	default:
// 		g.Cursor = false
// 		close(w.getDoneCh())
// 	}
// 	w.Mistyped = 0
// }

// func (w *Word) init(v *gocui.View) {
// 	w.e = newWordEditor()

// 	v.Editor = w.e
// 	v.Editable = true
// 	v.SelBgColor = gocui.ColorRed
// 	v.SelFgColor = gocui.ColorCyan
// }

// func newWordEditor() gocui.Editor {
// 	return gocui.EditorFunc(wordEditorFunc)
// }

// func wordEditorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
// 	word := strings.TrimSpace(getCurrentWord(v))

// 	switch {
// 	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
// 		handleDelete(v, true)
// 	case key == gocui.KeyDelete:
// 		handleDelete(v, false)
// 	case len(word) > maxWordLen:
// 		// do not add anymore runes
// 		// can only delete from here
// 	case ch != 0 && mod == 0:
// 		handleChar(v, ch)
// 	case key == gocui.KeySpace:
// 		handleSpace(v)
// 	}
// }

// func handleDelete(v *gocui.View, back bool) {
// 	v.EditDelete(back)
// 	checkAndHighlight(v)
// }

// func handleChar(v *gocui.View, ch rune) {
// 	v.EditWrite(ch)
// 	checkAndHighlight(v)
// }

// func checkAndHighlight(v *gocui.View) {
// 	w := strings.TrimSpace(getCurrentWord(v))
// 	ok := strings.HasPrefix(paragraph.CurrentWord(), w)
// 	highlight(ok, v)
// 	if !ok {
// 		word.Mistyped++
// 	}
// }

// func handleSpace(v *gocui.View) {
// 	w := strings.TrimSpace(getCurrentWord(v))
// 	if w == paragraph.CurrentWord() {
// 		clearEditor(v)

// 		perr := paragraph.Advance()
// 		if perr != nil {
// 			// finished typing all words by this points
// 			paragraph.Reset()
// 			word.Reset()
// 			stats.StopRace(true)
// 			controls.DefaultControlsContent()
// 			AfterRaceControls(g, false) // FIXME do not use global Gui ref
// 		}
// 	} else {
// 		// TODO Should we count mistyped space as mistake?
// 		highlight(false, v)
// 		v.EditWrite(' ')
// 	}
// }

// func clearEditor(v *gocui.View) {
// 	v.Clear()
// 	v.SetCursor(v.Origin())
// 	v.Editable = false
// }

// func highlight(ok bool, v *gocui.View) {
// 	v.Highlight = !ok
// 	paragraph.Mistyped = !ok
// }

// func getCurrentWord(v *gocui.View) string {
// 	line := v.Buffer()
// 	return line
// }
