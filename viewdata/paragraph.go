package viewdata

import (
	"errors"
	"strings"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/db"
	"go.uber.org/zap"
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

	// if set the wordeditor will be cleared
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
	return &ParagraphData{
		Words:    nil,
		wordi:    0,
		Mistyped: false,
	}
}

// StartRace is called when a race starts
func (pd *ParagraphData) StartRace(g *gocui.Gui, viewName string) error {
	if err := pd.setTargetParagraph(); err != nil {
		return err
	}

	pd.newDoneCh()
	pd.activateEditor(g, viewName)

	return nil
}

func (pd *ParagraphData) setTargetParagraph() error {
	para, err := db.ChooseParagraph()
	if err != nil {
		return err
	}

	pd.Words = strings.Fields(para)
	for i, w := range pd.Words {
		pd.Words[i] = strings.TrimSpace(w)
	}

	n := db.AddNewLines(pd.Words, pd.W-1)
	pd.lineCount = n
	pd.wordi = 0
	pd.Oy = 0
	pd.Line = 0
	pd.Word = 0

	return nil
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
	}

	return nil
}

// OnEditorChange is called on every change event to woreditor
func (pd *ParagraphData) OnEditorChange(w string) {
	endsWithSpace := strings.HasSuffix(w, " ")
	if endsWithSpace && len(w) > 1 {
		w = strings.TrimSuffix(w, " ")
	}

	cw := pd.currentWord()
	correct := strings.HasPrefix(cw, w)
	pd.Mistyped = !correct
	if endsWithSpace && w == cw {
		pd.tryAdvanceWord()
	}

	// Update UI and Stats
	pd.updateScrollAttr()
	pd.updateCh <- true
	pd.statsCh <- StatMsg{
		IsMistyped: pd.Mistyped,
	}
}

func (pd *ParagraphData) tryAdvanceWord() {
	if pd.wordi == len(pd.Words) {
		// TODO: end of race
		return
	}

	pd.wordi++
	pd.ShouldClearEditor = true
}

// DebugAdvance advances by a word to debug stuff
func (pd *ParagraphData) DebugAdvance() {
	pd.OnEditorChange(pd.currentWord() + " ")
}

// Called after Advancing by a word
func (pd *ParagraphData) updateScrollAttr() {
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
	prevWord := pd.Words[pd.wordi-1]
	config.Logger.Info("scroll attr", zap.Int("Word", pd.Word),
		zap.Int("Line", pd.Line), zap.Int("lineCount", pd.lineCount), zap.Int("Oy", pd.Oy), zap.String("prevWord", prevWord),
		zap.Bool("Line++", strings.HasSuffix(prevWord, "\n")))
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

func (pd *ParagraphData) activateEditor(g *gocui.Gui, viewName string) {
	g.SetCurrentView(viewName)
	g.Cursor = true
}

func (pd *ParagraphData) deactivateEditor(g *gocui.Gui) {
	g.Cursor = false
}
