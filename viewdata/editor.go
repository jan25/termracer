package viewdata

import (
	"errors"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	"go.uber.org/zap"
)

// MaxWordLen is maximum a word can go in editor view
const MaxWordLen int = 25

// WordEditorData stores content for word editor view
type WordEditorData struct {
	// ShouldClearEditor is set for every new target word
	ShouldClearEditor bool
	// use to check current word is mistyped
	IsMistyped bool

	// Channels to communicate with paragraph
	psender   chan WordValidateMsg
	preceiver chan WordValidateMsg

	// Channel to send messages to racestats
	rsender chan StatMsg

	done chan struct{}
}

// WordValidateMsg is used to communicate with ParagraphData
type WordValidateMsg struct {
	TypedWord  string
	Correct    bool
	IsNextWord bool
}

// NewWordEditorData creates new WordEditorData instance
func NewWordEditorData() *WordEditorData {
	return &WordEditorData{
		ShouldClearEditor: true,
		IsMistyped:        false,
	}
}

// StartRace starts a new race
func (w *WordEditorData) StartRace(g *gocui.Gui, viewName string) error {
	if w.psender == nil || w.preceiver == nil || w.rsender == nil {
		return errors.New("Channel for communication is nil")
	}

	w.newDoneCh()
	w.activateEditor(g, viewName)
	go w.talkWithParagraph()
	return nil
}

// FinishRace finishes a ongoing race
func (w *WordEditorData) FinishRace(g *gocui.Gui) error {
	w.ShouldClearEditor = true
	w.IsMistyped = false

	select {
	case <-w.DoneCh():
		return errors.New("Race already stopped")
	default:
		w.deactivateEditor(g)
		close(w.DoneCh())
	}
	return nil
}

// SetChannels sets channels for communication with other components during a race
func (w *WordEditorData) SetChannels(psender, preceiver chan WordValidateMsg, rsender chan StatMsg) {
	w.psender = psender
	w.preceiver = preceiver
	w.rsender = rsender
}

func (w *WordEditorData) talkWithParagraph() {
	defer close(w.preceiver) // FIXME: figure who closes what

	for {
		select {
		case <-w.DoneCh():
			return
		default:
			msg := <-w.preceiver
			w.understandMsg(msg)
			w.sendStatsUpdate(msg)
		}
	}
}

func (w *WordEditorData) understandMsg(msg WordValidateMsg) {
	w.IsMistyped = !msg.Correct
	w.ShouldClearEditor = msg.IsNextWord
}

func (w *WordEditorData) sendStatsUpdate(msg WordValidateMsg) {
	w.rsender <- StatMsg{
		IsMistyped: !msg.Correct,
	}
}

// OnChangeMsg sends a message to paragraph for onchange event
// TODO should this return a error for when psender is closed?
func (w *WordEditorData) OnChangeMsg(s string) {
	config.Logger.Info("OnChangeMsg", zap.String("s", s))
	w.psender <- WordValidateMsg{
		TypedWord: s,
	}
}

func (w *WordEditorData) newDoneCh() {
	w.done = make(chan struct{})
}

// DoneCh returns reference to done channel
func (w *WordEditorData) DoneCh() chan struct{} {
	if w.done == nil {
		w.newDoneCh()
	}
	return w.done
}

func (w *WordEditorData) activateEditor(g *gocui.Gui, viewName string) {
	g.SetCurrentView(viewName)
	g.Cursor = true
}

func (w *WordEditorData) deactivateEditor(g *gocui.Gui) {
	g.Cursor = false
}
