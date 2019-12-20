package main

import (
	"flag"
	"log"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/config"
	db "github.com/jan25/termracer/data"
)

func main() {
	// Flags
	debug := flag.Bool("debug", false, "flag for debug mode")
	flag.Parse()

	// Ensure the required data files on local FS are present
	if err := db.EnsureDataDirs(); err != nil {
		log.Panicln(err)
	}

	// Setup logger
	var err error
	_, err = config.InitLogger(*debug)
	if err != nil {
		log.Panicln(err)
	}
	defer config.Logger.Sync()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	app, err = initializeAppData(gui)
	if err != nil {
		log.Panicln(err)
	}

	// Setup key bindings on startup
	defaultBindings(gui)
	if *debug {
		debugBindings(gui, app)
	}

	config.Logger.Info("Starting main loop..")
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
