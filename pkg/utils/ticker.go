package utils

import (
	"time"

	"github.com/jan25/gocui"
)

// Tick starts the ticker
func Tick(t *Timer, forceUpdateChan chan bool, g *gocui.Gui) {
	t.Ticking = true

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer func() {
		t.Ticking = false
	}()
	defer close(forceUpdateChan)

	for {
		select {
		case <-t.getDoneCh():
			return
		case <-forceUpdateChan:
			UpdateUI(g)
		case <-ticker.C:
			UpdateUI(g)
		}
	}
}

// UpdateUI update the UI
// used when we want to manually/force update the UI
func UpdateUI(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		// do nothing
		// this automatically updates the timer in the view
		return nil
	})
}
