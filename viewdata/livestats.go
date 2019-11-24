package viewdata

import (
	"errors"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/pkg/utils"
)

// LiveStats represents contents show during a race
// in stats view
type LiveStats struct {
	// Count of In/Correctly typed chars
	correct   int
	incorrect int

	// stream of messages from word editor
	m chan StatMsg

	// done channel
	done chan struct{}

	IsActive bool // whether this data is active to show in stats view

	timer *utils.Timer
}

// StatMsg is used to communicate
// with wordeditor
type StatMsg struct {
	IsMistyped bool
}

// NewLiveStats creates new instance of LiveStats
func NewLiveStats(m *chan StatMsg) *LiveStats {
	return &LiveStats{
		correct:   0,
		incorrect: 0,
		m:         *m,
	}
}

// Start is called on new race start
func (ls *LiveStats) Start() error {
	if ls.m == nil {
		return errors.New("stream channel is nil")
	}

	if err := ls.timer.Start(); err != nil {
		return errors.New("Failed to start timer: " + err.Error())
	}

	go ls.listenToWordEditor()

	ls.IsActive = true
	ls.timer.Start()

	return nil
}

// TryStartTicker starts the ticker, so timer in stats view gets updated
// FIXME: anyway we can keep away from passing gocui instance?
func (ls *LiveStats) TryStartTicker(g *gocui.Gui) {
	if ls.timer.Ticking {
		return // do nothing
	}

	go utils.Tick(ls.timer, g)
}

func (ls *LiveStats) listenToWordEditor() {
	defer close(ls.m)

	for {
		select {
		case <-ls.getDoneCh():
			return
		default:
			msg := <-ls.m
			if msg.IsMistyped {
				ls.incorrect++
			} else {
				ls.correct++
			}
		}
	}
}

// Finish is called at end of a race
func (ls *LiveStats) Finish() error {
	var err error
	select {
	case <-ls.getDoneCh():
		return errors.New("race already stopped")
	default:
		close(ls.getDoneCh())
		err = ls.timer.Stop()
		if err != nil {
			return err
		}
		// TODO return current race stat
	}

	return nil
}

// ElapsedTime during a race
func (ls *LiveStats) ElapsedTime() (*utils.TimeFormatted, error) {
	return ls.timer.ElapsedTime()
}

// Wpm is words per minute stat from beginning of a race
func (ls *LiveStats) Wpm() int {
	return 0
}

// Accuracy of typing during a race
func (ls *LiveStats) Accuracy() float32 {
	return 0
}

func (ls *LiveStats) getDoneCh() chan struct{} {
	if ls.done == nil {
		ls.done = make(chan struct{})
	}
	return ls.done
}
