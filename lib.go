package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/jan25/termracer/server/client"
)

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

// AppendLineEOF appends a given line to end of a file
func AppendLineEOF(fname, line string) error {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		Logger.Error("Failed to append to file " + fname)
		return err
	}
	defer f.Close()

	// Add new line if not provided as part of line
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}

	if _, err := f.Write([]byte(line)); err != nil {
		Logger.Error("Failed to write to file " + fname)
		return err
	}
	return nil
}

// FormatDate formats time into
// DD-MM-YYYY format
func FormatDate(t time.Time) string {
	y, m, d := t.Date()
	// keep last 2 digits in year
	y %= 100
	return fmt.Sprintf("%02d/%02d/%d", d, m, y)
}

// CalculateWpm calculates words per minute
// based on words typed so far and time elapsed
func CalculateWpm(countw int, secs int) int {
	return (60 * paragraph.CountDoneWords()) / secs
}

// CalculateAccuracy calculates accuracy in a race
// based on chars typed so far and mistyped chars
func CalculateAccuracy(chars int, mistypes int) (float64, error) {
	if chars <= 0 {
		return 0, errors.New("Invalid chars param for CalculateAccuracy")
	}
	mistypesF := float64(mistypes)
	charsF := float64(chars)
	mistypesF = math.Min(mistypesF, charsF)
	return 100 * (1 - mistypesF/charsF), nil
}

// ChooseParagraph calls the server to fetch a paragraph
func ChooseParagraph() (string, error) {
	p, err := client.ChooseParagraph()
	if err != nil {
		Logger.Info("Error fetching paragraph from server " + err.Error())
		Logger.Info("Falling back to default paragraph")
		// TODO remove this temporary fix
		// returning const paragraph if couldn't reach server
		return firstParagraph, nil;
	}
	return p, nil
}

// GenerateLocalParagraphs checks if samples/use has > 0 paragraphs
// available. If not tries to generate them
func GenerateLocalParagraphs() error {
	files, err := ioutil.ReadDir("server/samples/use")
	if err != nil {
		return errors.New("failed to read use directory")
	}
	if len(files) == 0 {
		// TODO generate paragraphs if none available
	}
	return nil
}

const firstParagraph = `
She sank more and more into uneasy delirium. At times she shuddered,
turned her eyes from side to side, recognised everyone for a minute,
but at once sank into delirium again. Her breathing was hoarse and
difficult, there was a sort of rattle in her throat.
`
