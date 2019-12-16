package main

import (
	"log"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
)

// defaultBindings registers key bindings for default controls
func defaultBindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	afterRaceControls(g, true)
}

// duringRaceControls registers keybindings for a new race
func duringRaceControls(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone, ctrlE); err != nil {
		log.Panicln(err)
	}

	if err := g.DeleteKeybinding("", gocui.KeyCtrlS, gocui.ModNone); err != nil {
		log.Panicln(err)
	}
	if err := g.DeleteKeybinding("", gocui.KeyCtrlJ, gocui.ModNone); err != nil {
		log.Panicln(err)
	}
	if err := g.DeleteKeybinding("", gocui.KeyCtrlK, gocui.ModNone); err != nil {
		log.Panicln(err)
	}
}

// afterRaceControls registers keybindings at end of a race
func afterRaceControls(g *gocui.Gui, isDefault bool) {
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, ctrlS); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, app.HistoryScrollUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, app.HistoryScrollDown); err != nil {
		log.Panicln(err)
	}

	if !isDefault {
		if err := g.DeleteKeybinding("", gocui.KeyCtrlE, gocui.ModNone); err != nil {
			log.Panicln(err)
		}
	}
}

// handler for race start event
func ctrlS(g *gocui.Gui, v *gocui.View) error {
	err := app.OnRaceStart(g)
	duringRaceControls(g)
	return err
}

// handler for race end event
func ctrlE(g *gocui.Gui, v *gocui.View) error {
	err := app.OnRaceFinish()
	afterRaceControls(g, false)
	return err
}

// handler for quit app event
func quit(g *gocui.Gui, v *gocui.View) error {
	config.Logger.Info("Quitting termracer..")
	return gocui.ErrQuit
}

// this func will list keybindings for debugging purposes only
func debugBindings(g *gocui.Gui, ad *AppData) {
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, ad.DebugAdvance); err != nil {
		log.Panicln(err) // FIXME This might crash the app at end of paragraph
	}
}
