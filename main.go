package main

import (
	"log"
	"os"
	"path/filepath"

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

func initLogger() error {
	cfg := zap.NewProductionConfig()
	d, _ := GetTopLevelDir()
	path := d + "/logs/app.log"

	// create logs directory if not exists
	dirName := filepath.Dir(path)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}

	var _, err = os.Stat(path)
	// create log file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	cfg.OutputPaths = []string{path}
	Logger, _ = cfg.Build()
	return nil
}

// checks to see data dirs required for application are present
// creates dirs/files if not present
func ensureDataDirs() error {
	// ensure samples use directory
	s, err := GetSamplesUseDir()
	if err != nil {
		return err
	}
	if err := utils.CreateDirIfNotExists(s); err != nil {
		return err
	}
	if err := GenerateLocalParagraphs(); err != nil {
		return err
	}

	// ensure racehistory file
	rh, err := GetHistoryFilePath()
	if err != nil {
		return err
	}
	if err := utils.CreateFileIfNotExists(rh); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := ensureDataDirs(); err != nil {
		log.Panicln(err)
	}
	if err := initLogger(); err != nil {
		log.Panicln(err)
	}
	defer Logger.Sync()

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

	stats.InitKeyBindings(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, ctrlS); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone, ctrlE); err != nil {
		log.Panicln(err)
	}

	Logger.Info("Starting main loop...")
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func ctrlS(g *gocui.Gui, v *gocui.View) error {
	err := stats.StartRace()
	if err == nil {
		paragraph.Init()
		word.Init()
		controls.RaceModeControls()
	}
	// TODO ctrlS can be hit during the race
	// Handle the err from stats gracefully
	// Or disable these controls during race
	return nil
}

func ctrlE(g *gocui.Gui, v *gocui.View) error {
	paragraph.Reset()
	word.Reset()
	stats.StopRace(false)
	controls.DefaultControls()

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	Logger.Info("Quitting termracer..")
	return gocui.ErrQuit
}
