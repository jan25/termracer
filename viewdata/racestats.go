package viewdata

import "errors"

type StatsData struct {
	// Count of In/Correctly typed chars
	correct   int
	incorrect int

	// stream of messages from word editor
	m chan StatMsg

	// done channel
	done chan struct{}
}

type StatMsg struct {
	IsMistyped bool
}

func NewStatsData(m *chan StatMsg) *StatsData {
	return &StatsData{
		correct:   0,
		incorrect: 0,
		m:         *m,
	}
}

func (sd *StatsData) Start() error {
	if sd.m == nil {
		return errors.New("stream channel is nil")
	}

	go sd.listenToWordEditor()
	// TODO start timer

	return nil
}

func (sd *StatsData) Finish() error {
	select {
	case <-sd.getDoneCh():
		return errors.New("race already stopped")
	default:
		close(sd.getDoneCh())
		// TODO return current race stat
	}

	return nil
}

func (sd *StatsData) listenToWordEditor() {
	defer close(sd.m)

	for {
		select {
		case <-sd.getDoneCh():
			return
		default:
			msg := <-sd.m
			if msg.IsMistyped {
				sd.incorrect++
			} else {
				sd.correct++
			}
		}
	}
}

func (sd *StatsData) getDoneCh() chan struct{} {
	if sd.done == nil {
		sd.done = make(chan struct{})
	}
	return sd.done
}
