package main

import (
	"flag"
	"log"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/db"
	"github.com/jan25/termracer/pkg/utils"
	"go.uber.org/zap"
)

var (
	// Logger is a global file logger
	Logger *zap.Logger
	app    *AppData
)

func main() {
	// Flags
	debug := *flag.Bool("debug", false, "flag for debug mode")
	flag.Parse()

	// Setup logger
	var err error
	Logger, err = utils.InitLogger("./app.log", debug)
	if err != nil {
		log.Panicln(err)
	}
	defer Logger.Sync()

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

	if debug {
		debugBindings(gui)
	}

	Logger.Info("Starting main loop..")
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
