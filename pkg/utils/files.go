package utils

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
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
func CreateFileIfNotExists(fpath string) error {
	var _, err = os.Stat(fpath)
	// create log file if not exists
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(fpath)
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

// AppendLineEOF appends a given line to end of a file
func AppendLineEOF(fname, line string) error {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Add new line if not provided as part of line
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}

	if _, err := f.Write([]byte(line)); err != nil {
		return err
	}
	return nil
}

// WriteToFile writes to a file at a given path
func WriteToFile(fpath string, data []byte) error {
	err := ioutil.WriteFile(fpath, data, 0644)
	return err
}
