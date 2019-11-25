package viewdata

// MaxWordLen is maximum a word can go in editor view
const MaxWordLen int = 25

// WordEditorData stores content for word editor view
type WordEditorData struct {
	// CurrentTyped is the word currently typed in editor
	CurrentTyped string
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
	CurrentTyped string
	Correct      bool
	IsNewWord    bool
}

// NewWordEditorData creates new WordEditorData instance
func NewWordEditorData() *WordEditorData {
	return &WordEditorData{
		CurrentTyped: "",
		IsMistyped:   false,
	}
}

// StartRace starts a new race
func (w *WordEditorData) StartRace(psender, preceiver chan WordValidateMsg, rsender chan StatMsg) {
	w.psender = psender
	w.preceiver = preceiver
	w.rsender = rsender

	w.newDoneCh()

	go w.talkWithParagraph()
}

// FinishRace finishes a ongoing race
func (w *WordEditorData) FinishRace() {
	w.CurrentTyped = ""
	w.IsMistyped = false

	close(w.DoneCh())
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
	w.IsMistyped = msg.Correct
	w.CurrentTyped = msg.CurrentTyped
}

func (w *WordEditorData) sendStatsUpdate(msg WordValidateMsg) {
	w.rsender <- StatMsg{
		IsMistyped: msg.Correct,
	}
}

// OnChangeMsg sends a message to paragraph for onchange event
// TODO should this return a error for when psender is closed?
func (w *WordEditorData) OnChangeMsg(s string) {
	w.psender <- WordValidateMsg{
		CurrentTyped: s,
		Correct:      true,
		IsNewWord:    false,
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
