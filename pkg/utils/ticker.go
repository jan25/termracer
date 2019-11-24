package utils

import (
	"time"

	"github.com/jan25/gocui"
)

// Tick starts the ticker
func Tick(t *Timer, g *gocui.Gui) {
	t.Ticking = true

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer func() {
		t.Ticking = false
	}()

	for {
		select {
		case <-t.getDoneCh():
			return
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				// do nothing
				// this automatically updates the timer in the view
				return nil
			})
		}
	}
}
