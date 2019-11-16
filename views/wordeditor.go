package views

import (
	"github.com/jan25/gocui"
	"github.com/jan25/termracer/viewdata"
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
	Data *viewdata.WordEditorData
}

func newWordView(name string, x, y int, w, h int) *WordView {
	wv := &WordView{
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
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
	case <-w.getDoneCh():
		// channel closed
		w.clearEditor(v)
	default:
		w.initEditor(v, &w.e)
	}
	return nil
}

func (w *WordView) getDoneCh() chan struct{} {
	if w.done == nil {
		w.done = make(chan struct{})
	}
	return w.done
}

func (w *WordView) initEditor(v *gocui.View, e *gocui.Editor) {
	v.Editor = *e
	v.Editable = true
	v.SelBgColor = gocui.ColorRed
	v.SelFgColor = gocui.ColorCyan

	if len(w.Data.CurrentTyped) == 0 {
		w.clearEditor(v) // Reset to origin for new target word
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
	case len(currentTyped) > viewdata.MaxWordLen:
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
	w.onChange(v)

	// TODO figure how to finish race at end of target paragraph
}

// sends message to paragraph for
// typed word validations
func (w *WordView) onChange(v *gocui.View) {
	w.Data.OnChangeMsg(w.getCurrentTyped(v))
}

// highlight for incorrectly typed word
func (w *WordView) highlight(v *gocui.View) {
	v.Highlight = w.Data.IsMistyped
}

// gets current word in the editor
func (w *WordView) getCurrentTyped(v *gocui.View) string {
	line := v.Buffer()
	return line
}
