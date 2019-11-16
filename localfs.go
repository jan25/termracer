package main

import "github.com/jan25/termracer/pkg/utils"

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
