package viewdata

import (
	"errors"
	"strings"

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

	// Channels to communicate with wordeditor
	wsender   chan WordValidateMsg
	wreceiver chan WordValidateMsg

	// channel to update UI
	updateCh chan bool

	done chan struct{}
}

// NewParagraphData creates instance of ParagraphData
func NewParagraphData() *ParagraphData {
	return &ParagraphData{
		Words:     nil,
		wordi:     0,
		Mistyped:  false,
		lineCount: 0, // FIXME
		Line:      0,
		Word:      0,
		Oy:        0,
	}
}

// StartRace is called when a race starts
func (pd *ParagraphData) StartRace() error {
	if pd.wsender == nil || pd.wreceiver == nil {
		return errors.New("wsender or wreceiver is nil")
	}

	if err := pd.setTargetParagraph(); err != nil {
		return err
	}

	pd.newDoneCh()

	go pd.talkWithWordEditor()

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

// FinishRace is called to finish a race
func (pd *ParagraphData) FinishRace() error {
	select {
	case <-pd.DoneCh():
		return errors.New("race already stopped")
	default:
		close(pd.DoneCh())
	}

	return nil
}

// SetChannels sets channels for communication
func (pd *ParagraphData) SetChannels(wsender, wreceiver chan WordValidateMsg, updateCh chan bool) {
	pd.wsender = wsender
	pd.wreceiver = wreceiver
	pd.updateCh = updateCh
}

func (pd *ParagraphData) talkWithWordEditor() {
	defer close(pd.wreceiver)

	for {
		select {
		case <-pd.DoneCh():
			return
		default:
			msg := <-pd.wreceiver
			pd.validateTypedWord(msg)
		}
	}
}

func (pd *ParagraphData) validateTypedWord(msg WordValidateMsg) {
	s := strings.TrimSuffix(msg.TypedWord, " ") // trim single space in suffix
	cw := pd.currentWord()

	correct := strings.HasPrefix(cw, s)
	newWord := (s == cw) && strings.HasSuffix(msg.TypedWord, " ")

	pd.wsender <- WordValidateMsg{
		Correct:    correct,
		IsNextWord: newWord,
	}

	pd.Mistyped = !correct
	config.Logger.Info("setting pd.Mistyped", zap.Bool("mistyped", pd.Mistyped), zap.String("s", s), zap.String("cw", cw))

	if newWord {
		pd.tryAdvanceWord()
	}
}

func (pd *ParagraphData) tryAdvanceWord() {
	config.Logger.Info("Advancing target word")
	if pd.wordi == len(pd.Words) {
		// TODO: end of race
		return
	}

	pd.wordi++
	pd.updateCh <- true
}

func (pd *ParagraphData) currentWord() string {
	w := pd.Words[pd.wordi]
	w = strings.TrimSuffix(w, "\n")
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
