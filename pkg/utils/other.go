package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// FormatDate formats time into
// DD-MM-YYYY format
func FormatDate(t time.Time) string {
	y, m, d := t.Date()
	// keep last 2 digits in year
	y %= 100
	return fmt.Sprintf("%02d/%02d/%d", d, m, y)
}

// InitLogger initalizes a logger
func InitLogger(fpath string, debug bool) (*zap.Logger, error) {
	if !debug {
		// no-op logger
		return zap.New(nil), nil
	}

	cfg := zap.NewProductionConfig()
	if err := CreateFileIfNotExists(fpath); err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{fpath}
	l, _ := cfg.Build()
	return l, nil
}
