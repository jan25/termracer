package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

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
	randf := files[rand.Int31n(int32(n))]

	b, err := ioutil.ReadFile("samples/use/" + randf.Name())
	if err != nil {
		return "ERROR", errors.New("error in reading a paragraph file")
	}
	return string(b), nil
}
