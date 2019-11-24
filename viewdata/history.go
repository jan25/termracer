package viewdata

import "time"

// OneStat represents information
// about completed race
type OneStat struct {
	// words per minute
	Wpm int
	// accuracy eg. 95.6%
	Accuracy float64
	// time of race
	When time.Time
}

// Stats is a datastructure to
// store stats for past races
type Stats struct {
	History  []*OneStat
	Selected int // To keep track of highlighted stat in history
}
