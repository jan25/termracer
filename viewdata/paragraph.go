package viewdata

import (
	"errors"
	"strings"
)

// ParagraphData keeps track of state, data
// required for view's content
type ParagraphData struct {
	// words in paragraph
	words []string
	// index of target word - 0 based
	wordi int
	// whether target word is mistyped
	mistyped bool

	// line count in target paragraph
	lineCount int

	// For highlighted word; Line, Word number both start at 0
	Line int
	Word int
}

// Words returns the array of words in target paragraph
func (pd *ParagraphData) Words() []string {
	return pd.words
}

// TargetWordIndex returns index of target word in paragraph
func (pd *ParagraphData) TargetWordIndex() int {
	return pd.wordi
}

// IsMistyped tells if target word is being mistyped
func (pd *ParagraphData) IsMistyped() bool {
	return pd.mistyped
}

// LineCount returns count of lines shown in the view
func (pd *ParagraphData) LineCount() int {
	return pd.lineCount
}

// Advance moves target word to next word
func (pd *ParagraphData) Advance() error {
	if pd.wordi >= len(pd.words)-1 {
		return errors.New("can not advance beyond number of words")
	}

	pd.wordi++

	// Scroll logic
	prevWord := pd.words[pd.wordi-1]
	if strings.HasSuffix(prevWord, "\n") {
		// Jump to next line
		pd.Line++
		pd.Word = 0
	} else {
		pd.Word++
	}

	// FIXME Put below this inside views/paragraph
	// p.makeScroll()

	return nil
}
