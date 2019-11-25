package viewdata

import (
	"errors"
	"strings"
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

	// For highlighted word; Line, Word number both start at 0
	Line int
	Word int

	// Channels to communicate with wordeditor
	wsender   chan WordValidateMsg
	wreceiver chan WordValidateMsg

	done chan struct{}
}

// NewParagraphData creates instance of ParagraphData
func NewParagraphData(sender, receiver *chan WordValidateMsg) *ParagraphData {
	return &ParagraphData{
		Words:     getTargetWords(),
		wordi:     0,
		Mistyped:  false,
		lineCount: 0, // FIXME
		Line:      0,
		Word:      0,
		wsender:   *sender,
		wreceiver: *receiver,
	}
}

func getTargetWords() []string {
	// TODO choose paragraph and do some work
	return nil
}

// Start is called when a race starts
func (pd *ParagraphData) Start() error {
	if pd.wsender == nil || pd.wreceiver == nil {
		return errors.New("wsender or wreceiver is nil")
	}

	pd.newDoneCh()

	go pd.talkWithWordEditor()

	return nil
}

// Finish is called to finish a race
func (pd *ParagraphData) Finish() error {
	select {
	case <-pd.DoneCh():
		return errors.New("race already stopped")
	default:
		close(pd.DoneCh())
	}

	return nil
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
	s := strings.TrimSuffix(msg.CurrentTyped, " ") // time any suffix at end
	cw := pd.currentWord()

	correct := strings.HasPrefix(s, cw)
	newWord := (s == cw) && strings.HasSuffix(msg.CurrentTyped, " ")
	setTyped := msg.CurrentTyped
	if newWord {
		setTyped = "" // resets editor
	}

	pd.wsender <- WordValidateMsg{
		Correct:      correct,
		IsNewWord:    newWord,
		CurrentTyped: setTyped,
	}

	pd.Mistyped = correct
}

func (pd *ParagraphData) currentWord() string {
	return pd.Words[pd.wordi]
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
