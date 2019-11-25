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
	wreceiver chan StatMsg

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
func NewLiveStats() *LiveStats {
	return &LiveStats{
		correct:   0,
		incorrect: 0,
		IsActive:  false, // default: no race in progress at app startup
	}
}

// StartRace starts a new race
func (ls *LiveStats) StartRace() error {
	if ls.wreceiver == nil {
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
	defer close(ls.wreceiver)

	for {
		select {
		case <-ls.getDoneCh():
			return
		default:
			msg := <-ls.wreceiver
			if msg.IsMistyped {
				ls.incorrect++
			} else {
				ls.correct++
			}
		}
	}
}

// FinishRace finishes a ongoing race
func (ls *LiveStats) FinishRace() error {
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
func (ls *LiveStats) Wpm() (int, error) {
	secs, err := ls.timer.ElapsedTimeInSecs()
	if err != nil {
		return 0, err
	}
	wpm := utils.CalculateWpm(ls.correct, secs)
	return int(wpm), nil
}

// Accuracy of typing during a race
func (ls *LiveStats) Accuracy() (float64, error) {
	return utils.CalculateAccuracy(ls.correct+ls.incorrect, ls.incorrect)
}

func (ls *LiveStats) getDoneCh() chan struct{} {
	if ls.done == nil {
		ls.done = make(chan struct{})
	}
	return ls.done
}
