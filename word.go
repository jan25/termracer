package main

import (
	"strings"

	"github.com/jan25/gocui"
)

const MAX_WORD_LEN int = 15

// for developerment
// remove when this module is fully developed
const TEST_WORD string = "testword"

var WordEditor gocui.Editor = gocui.EditorFunc(wordEditorFunc)

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
		}
		paragraph.DrawView()
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
