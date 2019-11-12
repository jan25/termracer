package viewdata

// MaxWordLen is maximum a word can go in editor view
const MaxWordLen int = 15

// WordEditorData stores content for word editor view
type WordEditorData struct {
	// use to check current word is mistyped
	isMistyped bool

	// Channels to communicate with paragraph
	psender   chan WordValidateMsg
	preceiver chan WordValidateMsg

	// Channel to send messages to racestats
	rsender chan StatMsg

	done chan struct{}
}

// WordValidateMsg is used to communicate with ParagraphData
type WordValidateMsg struct {
	Current   string
	Correct   bool
	IsNewWord bool
}

func (w *WordEditorData) talkWithParagraph() {
	defer close(w.preceiver)

	for {
		select {
		case <-w.getDoneCh():
			return
		default:
			msg := <-w.preceiver
			w.understandMsg(msg)
		}
	}
}

func (w *WordEditorData) understandMsg(msg WordValidateMsg) {
	// TODO
}

func (w *WordEditorData) getDoneCh() chan struct{} {
	if w.done == nil {
		w.done = make(chan struct{})
	}
	return w.done
}
