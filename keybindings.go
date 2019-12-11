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

	AfterRaceControls(g, true)
}

// DuringRaceControls registers keybindings when a race is in progress
func DuringRaceControls(g *gocui.Gui) {
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

// AfterRaceControls registers keybindings when no race in progress
func AfterRaceControls(g *gocui.Gui, isDefault bool) {
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

func ctrlS(g *gocui.Gui, v *gocui.View) error {
	err := app.OnRaceStart(g)
	DuringRaceControls(g)
	return err
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	err := app.OnRaceFinish()
	AfterRaceControls(g, false)
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	config.Logger.Info("Quitting termracer..")
	return gocui.ErrQuit
}

func debugBindings(g *gocui.Gui, ad *AppData) {
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, ad.DebugAdvance); err != nil {
		log.Panicln(err) // FIXME This might crash the app at end of paragraph
	}
}
