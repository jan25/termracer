package main

import (
	"flag"
	"log"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/db"
)

var (
	app *AppData
)

func main() {
	// Flags
	debug := flag.Bool("debug", false, "flag for debug mode")
	flag.Parse()

	// Setup logger
	var err error
	_, err = config.InitLogger(*debug)
	if err != nil {
		log.Panicln(err)
	}
	defer config.Logger.Sync()

	// Ensure the required data files on local FS are present
	if err := db.EnsureDataDirs(); err != nil {
		log.Panicln(err)
	}

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	app, err = InitializeAppData(gui)
	if err != nil {
		log.Panicln(err)
	}

	// Default key bindings on startup
	DefaultBindings(gui)

	if *debug {
		debugBindings(gui)
	}

	config.Logger.Info("Starting main loop..")
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
