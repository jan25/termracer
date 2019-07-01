package main

import "os/user"

const (
	// TopLevelDir is name of directory that stores data for termracer
	TopLevelDir = "/termracer"

	// SamplesDir is a sub directory to store paragraph samples data
	SamplesDir = "/samples/use"

	// HistoryFile is name of file to store race history data
	HistoryFile = "/racehistory.csv"
)

// GetSamplesUseDir returns path to samples/use
// directory in local filesystem
func GetSamplesUseDir() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir + SamplesDir, nil
}

// GetHistoryFilePath returns path to race history file
func GetHistoryFilePath() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	return home + TopLevelDir + HistoryFile, nil
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
