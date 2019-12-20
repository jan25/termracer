package data

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/jan25/gocui"
	db "github.com/jan25/termracer/data"
	"github.com/jan25/termracer/pkg/wordwrap"
)

// ParagraphData keeps track of state, data
// required for view's content
type ParagraphData struct {
	// words in paragraph
	Words []string
	// index of target word - 0 based
	wordi int
	// whether target word is mistyped
	Mistyped bool

	// true if a race is in progress
	RaceInProgress bool
	// if set the wordeditor will be cleared for next target word
	ShouldClearEditor bool

	// line count in target paragraph
	lineCount int

	// Dimensions of the view
	// FIXME: this shouldn't be available in viewdata
	H int
	W int
	// For highlighted word; Line, Word number both start at 0
	Line int
	Word int
	// Y position of View origin
	Oy int

	// channel to update UI
	updateCh chan bool

	// channel for sending LiveStats data
	statsCh chan StatMsg

	done chan struct{}
}

// NewParagraphData creates instance of ParagraphData
func NewParagraphData() *ParagraphData {
	// FIXME improve algorithm to randomise chosen paragraph
	// this seed gives no good random behavior
	// noticed repeated paragraphs few times
	rand.Seed(time.Now().Unix())

	return &ParagraphData{
		Words:          nil,
		wordi:          0,
		Mistyped:       false,
		RaceInProgress: false,
	}
}

// StartRace is called when a race starts
func (pd *ParagraphData) StartRace(g *gocui.Gui) {
	pd.setTargetParagraph()
	pd.newDoneCh()
	pd.RaceInProgress = true
}

func (pd *ParagraphData) setTargetParagraph() {
	para := db.ChooseParagraph()

	pd.Words = strings.Fields(para)
	for i, w := range pd.Words {
		pd.Words[i] = strings.TrimSpace(w)
	}

	n := addNewLines(pd.Words, pd.W-1)
	pd.lineCount = n
	pd.wordi = 0
	pd.Oy = 0
	pd.Line = 0
	pd.Word = 0
}

// SetChannels sets channels for UI, stats updates
func (pd *ParagraphData) SetChannels(statsCh chan StatMsg, updateCh chan bool) {
	pd.updateCh = updateCh
	pd.statsCh = statsCh
}

// FinishRace is called to finish a race
func (pd *ParagraphData) FinishRace() error {
	select {
	case <-pd.DoneCh():
		return errors.New("Race already stopped")
	default:
		close(pd.DoneCh())
		pd.Mistyped = false
		pd.updateCh <- true // this tick updates after race finish
	}

	return nil
}

// OnEditorChange is called on every change event to woreditor
func (pd *ParagraphData) OnEditorChange(w string) {
	endsWithSpace := strings.HasSuffix(w, " ")
	if endsWithSpace && len(w) > 1 {
		w = strings.TrimSuffix(w, " ")
	}

	lastWord := false
	if pd.wordi == len(pd.Words)-1 {
		lastWord = true
	}

	cw := pd.currentWord()
	correct := strings.HasPrefix(cw, w)
	pd.Mistyped = !correct
	if (lastWord || endsWithSpace) && w == cw {
		if pd.tryAdvanceWord() == false {
			return // end of race
		}
	}

	// Update UI and Stats
	pd.updateScrollAttr()
	pd.updateCh <- true
	pd.statsCh <- StatMsg{
		IsMistyped: pd.Mistyped,
		FinishRace: false,
	}
}

func (pd *ParagraphData) tryAdvanceWord() bool {
	if pd.wordi+1 == len(pd.Words) {
		// End of race, we ran out of words to type
		pd.statsCh <- StatMsg{
			FinishRace: true,
		}
		return false
	}

	pd.wordi++
	pd.ShouldClearEditor = true
	return true
}

// DebugAdvance advances by a word to debug stuff
func (pd *ParagraphData) DebugAdvance() {
	pd.OnEditorChange(pd.currentWord() + " ")
}

// Called after Advancing by a word
func (pd *ParagraphData) updateScrollAttr() {
	if pd.wordi <= 0 {
		return
	}
	prevWord := pd.Words[pd.wordi-1]
	if strings.HasSuffix(prevWord, "\n") {
		pd.Word = 0
		pd.Line++
	} else {
		pd.Word++
	}
	pd.makeScroll()
}

func (pd *ParagraphData) currentWord() string {
	w := pd.Words[pd.wordi]
	w = strings.TrimSuffix(w, "\n") // FIXME: do we need \n at end of a word? maybe this is to print on a new line
	return w
}

// GetCurrentIdx tells index of target word
func (pd *ParagraphData) GetCurrentIdx() int {
	return pd.wordi
}

// GetLineCount returns number of lines displayed in
// paragraph view
func (pd *ParagraphData) GetLineCount() int {
	return pd.lineCount
}

func (pd *ParagraphData) newDoneCh() {
	pd.done = make(chan struct{})
}

// DoneCh returns reference to done channel
func (pd *ParagraphData) DoneCh() chan struct{} {
	if pd.done == nil {
		pd.newDoneCh()
	}
	return pd.done
}

// makeScroll auto scrolls the view during the race
func (pd *ParagraphData) makeScroll() {
	whenLinesLeft := 2
	atWord, atLine := 2, (pd.H-1)-whenLinesLeft

	if pd.Word != atWord {
		return
	}
	currLine := pd.Line - pd.Oy
	linesLeft := (pd.GetLineCount() - 1) - pd.Line
	if currLine == atLine && linesLeft >= whenLinesLeft {
		pd.Oy++
	}
}

// addNewLines adds new line char to certian words
// to wrap and align the words into seperate lines
func addNewLines(words []string, width int) int {
	processed := []wordwrap.Word{}
	for _, w := range words {
		processed = append(processed, wordwrap.Word{
			Len: len(w),
		})
	}

	wordwrap.Wrap(processed, width)
	lines := 1
	for i, w := range processed {
		if w.Wrap {
			words[i] = words[i] + "\n"
			lines++
		}
	}
	return lines
}
