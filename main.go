package main

import (
	"flag"
	"log"

	"github.com/jan25/gocui"
	"github.com/jan25/termracer/pkg/utils"
	"go.uber.org/zap"
)

const (
	statsName    = "stats"
	paraName     = "para"
	wordName     = "word"
	controlsName = "controls"
)

var (
	// Logger is a global file logger
	Logger    *zap.Logger
	g         *gocui.Gui
	paragraph *Paragraph
	word      *Word
	stats     *StatsView
	controls  *Controls
)

var (
	paraW, paraH = 60, 8
	wordW, wordH = 60, 2

	statsW, statsH       = 20, 6
	controlsW, controlsH = 20, 4

	topX, topY = 1, 1
	pad        = 1
)

func main() {
	// Flags
	flag.Parse()
	debug := *flag.Bool("debug", false, "flag for debug mode")

	// Setup logger
	var err error
	Logger, err = utils.InitLogger("./app.log", debug)
	if err != nil {
		log.Panicln(err)
	}
	defer Logger.Sync()

	// Ensure the required data files on local FS are present
	if err := ensureDataDirs(); err != nil {
		log.Panicln(err)
	}

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g = gui
	defer g.Close()

	paragraph = newParagraph(paraName, topX, topY, paraW, paraH)
	word = newWord(wordName, topX, topY+paraH+pad, wordW, wordH)
	stats = newStatsView(statsName, topX+paraW+pad, topY, statsW, statsH)
	controls = newControls(controlsName, topX+paraW+pad, topY+statsH+pad, controlsW, controlsH)

	g.SetManager(paragraph, word, stats, controls)

	// Default key bindings on startup
	DefaultBindings(gui)

	if debug {
		debugBindings(gui)
	}

	Logger.Info("Starting main loop..")
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
