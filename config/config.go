package config

import "github.com/jan25/termracer/pkg/utils"

// TODO: Use pkg/path from stdlib instead of strings
const (
	// TopLevelDir is name of directory that stores data for termracer
	TopLevelDir = "/termracer"

	// SamplesJSONPath is a file storing sample paragraphs in JSON format
	SamplesJSONPath = "/samples.json"

	// HistoryFile is name of file to store race history data
	HistoryFile = "/racehistory.csv"

	// DebugLogs is path for logs for debugging purposes
	DebugLogs = "/app.log"

	// MaxWordLen is assumed/allowed maximum length of typed word
	MaxWordLen = 25
)

// GetTopLevelDir returns full path to top level dir
func GetTopLevelDir() (string, error) {
	home, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir, nil
}

// GetSamplesFilePath returns absolute path to samples.json file
func GetSamplesFilePath() (string, error) {
	tld, err := GetTopLevelDir()
	if err != nil {
		return "", err
	}
	return tld + SamplesJSONPath, nil
}

// GetHistoryFilePath returns path to race history file
func GetHistoryFilePath() (string, error) {
	tld, err := GetTopLevelDir()
	if err != nil {
		return "", err
	}
	return tld + HistoryFile, nil
}

// GetLogsFilePath is absolute path to log file
func GetLogsFilePath() (string, error) {
	tld, err := GetTopLevelDir()
	if err != nil {
		return "", err
	}
	return tld + DebugLogs, nil
}
