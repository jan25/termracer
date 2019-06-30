package main

import (
	"strings"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

// CreateFileIfNotExists creates file if not exists
func CreateFileIfNotExists(fname string) error {
	var _, err = os.Stat(fname)
	// create log file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer file.Close()
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

// ChooseParagraph chooses a paragraph from available
// paragraphs under /samples/use. This paragraph will
// be used in the next race and shown in paragraph View
func ChooseParagraph() (string, error) {
	files, err := ioutil.ReadDir("samples/use")
	n := len(files)
	if err != nil {
		return "ERROR", errors.New("failed to read use directory")
	}

	rand.Seed(time.Now().Unix())
	randf := files[rand.Int31n(int32(n))]

	b, err := ioutil.ReadFile("samples/use/" + randf.Name())
	if err != nil {
		return "ERROR", errors.New("error in reading a paragraph file")
	}
	return string(b), nil
}
