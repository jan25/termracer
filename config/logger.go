package config

import (
	"github.com/jan25/termracer/pkg/utils"
	"go.uber.org/zap"
)

// Logger is global logger instance
var Logger *zap.Logger

// InitLogger initalizes a logger
func InitLogger(fpath string, debug bool) (*zap.Logger, error) {
	if !debug {
		// no-op logger
		Logger = zap.New(nil)
		return Logger, nil
	}

	cfg := zap.NewProductionConfig()
	if err := utils.CreateFileIfNotExists(fpath); err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{fpath}
	Logger, _ := cfg.Build()
	return Logger, nil
}
