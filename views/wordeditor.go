package views

import (
	"strings"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	viewdata "github.com/jan25/termracer/views/data"
)

// WordView is a editor widget
type WordView struct {
	// name of the View
	name string
	// position, dimentions
	x, y int
	w, h int

	// editor instance
	e gocui.Editor

	// keep track of status of view
	done chan struct{}

	// data about contents of editor view
	Data *viewdata.ParagraphData
}

// NewWordView creates new instance of WordView
func NewWordView(name string, x, y int, w, h int, pData *viewdata.ParagraphData) *WordView {
	wv := &WordView{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
		Data: pData,
	}
	wv.e = wv.newWordEditor() // This looks wierd, doesn't it?
	return wv
}

// Layout is layout manager for word widget
func (w *WordView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	select {
	case <-w.Data.DoneCh():
		// no race in progress
		w.clearEditor(v)
		w.deactivateEditor(g)
	default:
		w.initEditor(v, &w.e, g)
	}
	return nil
}

func (w *WordView) initEditor(v *gocui.View, e *gocui.Editor, g *gocui.Gui) {
	if !w.Data.RaceInProgress {
		// this is on app startup
		w.deactivateEditor(g)
		return
	}

	w.activateEditor(g, w.name)
	v.Editor = *e
	v.Editable = true
	v.SelBgColor = gocui.ColorRed
	v.SelFgColor = gocui.ColorCyan

	if w.Data.ShouldClearEditor {
		w.clearEditor(v)                 // Reset to origin for new target word
		w.Data.ShouldClearEditor = false // Reset only once. This removes the deadlock on editor
	} else {
		w.highlight(v)
	}
}

func (w *WordView) clearEditor(v *gocui.View) {
	v.Clear()
	v.SetCursor(v.Origin())
	v.Editable = false
}

func (w *WordView) newWordEditor() gocui.Editor {
	return gocui.EditorFunc(w.wordEditorFunc)
}

func (w *WordView) wordEditorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	currentTyped := w.getCurrentTyped(v)

	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		w.handleDelete(v, true)
	case key == gocui.KeyDelete:
		w.handleDelete(v, false) // FIXME figure why did we write this?
	case len(currentTyped) > config.MaxWordLen:
		// do not add anymore runes
		// can only delete from here
	case ch != 0 && mod == 0:
		w.handleChar(v, ch)
	case key == gocui.KeySpace:
		w.handleSpace(v)
	}
}

func (w *WordView) handleDelete(v *gocui.View, back bool) {
	v.EditDelete(back)
	w.onChange(v)
}

func (w *WordView) handleChar(v *gocui.View, ch rune) {
	v.EditWrite(ch)
	w.onChange(v)
}

func (w *WordView) handleSpace(v *gocui.View) {
	v.EditWrite(' ') // single space
	w.onChange(v)
}

// sends message to paragraph for
// typed word validations
func (w *WordView) onChange(v *gocui.View) {
	w.Data.OnEditorChange(w.getCurrentTyped(v))
}

// highlight for incorrectly typed word
func (w *WordView) highlight(v *gocui.View) {
	v.Highlight = w.Data.Mistyped
}

// gets current word in the editor
func (w *WordView) getCurrentTyped(v *gocui.View) string {
	line := v.Buffer()
	line = strings.TrimSuffix(line, "\n") // remove new line thingy at end of buffer
	return line
}

func (w *WordView) activateEditor(g *gocui.Gui, viewName string) {
	g.SetCurrentView(viewName)
	g.Cursor = true
}

// deactivateEditor deactivated editor view
func (w *WordView) deactivateEditor(g *gocui.Gui) {
	g.Cursor = false
}
