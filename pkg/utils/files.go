package utils

import (
	"os"
	"os/user"
)

// GetHomeDir returns absolute path to Home directory
func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// CreateFileIfNotExists creates file if not exists
func CreateFileIfNotExists(fname string) error {
	var _, err = os.Stat(fname)
	// create log file if not exists
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// CreateDirIfNotExists creates directory path if not exists
func CreateDirIfNotExists(dpath string) error {
	var _, err = os.Stat(dpath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
