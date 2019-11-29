package config

import "github.com/jan25/termracer/pkg/utils"

const (
	// TopLevelDir is name of directory that stores data for termracer
	TopLevelDir = "/termracer"

	// SamplesDir is a sub directory to store paragraph samples data
	SamplesDir = "/samples/use"

	// SamplesJSONPath is a file storing sample paragraphs in JSON format
	SamplesJSONPath = SamplesDir + "/samples.json"

	// HistoryFile is name of file to store race history data
	HistoryFile = "/racehistory.csv"

	// DebugLogs is path for logs for debugging purposes
	DebugLogs = "/app.log"
)

// GetSamplesUseDir returns path to samples/use
// directory in local filesystem
func GetSamplesUseDir() (string, error) {
	home, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir + SamplesDir, nil
}

// GetSamplesFilePath returns absolute path to samples.json file
func GetSamplesFilePath() (string, error) {
	home, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir + SamplesJSONPath, nil
}

// GetHistoryFilePath returns path to race history file
func GetHistoryFilePath() (string, error) {
	home, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir + HistoryFile, nil
}

// GetTopLevelDir returns full path to top level dir
func GetTopLevelDir() (string, error) {
	home, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir, nil
}

// GetLogsFilePath is absolute path to log file
func GetLogsFilePath() (string, error) {
	tld, err := GetTopLevelDir()
	if err != nil {
		return "", err
	}
	return tld + DebugLogs, nil
}
