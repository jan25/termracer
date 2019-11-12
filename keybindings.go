package main

import (
	"log"

	"github.com/jan25/gocui"
)

// DefaultBindings registers key bindings for default controls
func DefaultBindings(g *gocui.Gui) {
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
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, stats.ScrollUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, stats.ScrollDown); err != nil {
		log.Panicln(err)
	}

	if !isDefault {
		if err := g.DeleteKeybinding("", gocui.KeyCtrlE, gocui.ModNone); err != nil {
			log.Panicln(err)
		}
	}
}

func ctrlS(g *gocui.Gui, v *gocui.View) error {
	DuringRaceControls(g)
	stats.StartRace()
	paragraph.Init()
	word.Init()
	controls.RaceModeControlsContent()

	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	AfterRaceControls(g, false)
	paragraph.Reset()
	word.Reset()
	stats.StopRace(false)
	controls.DefaultControlsContent()

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	Logger.Info("Quitting termracer..")
	return gocui.ErrQuit
}

func debugBindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, advanceWord); err != nil {
		log.Panicln(err) // FIXME This will crash the app at end of paragraph
	}
}

// For debugging in the UI
func advanceWord(g *gocui.Gui, v *gocui.View) error {
	return paragraph.Advance()
}
