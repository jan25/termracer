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
		return zap.New(nil), nil
	}

	cfg := zap.NewProductionConfig()
	if err := utils.CreateFileIfNotExists(fpath); err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{fpath}
	l, _ := cfg.Build()
	return l, nil
}
