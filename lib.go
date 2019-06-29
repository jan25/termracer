package main

import (
	"errors"
	"io/ioutil"
	"math"
)

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
	// TODO choose random paragraph from use directory
	b, err := ioutil.ReadFile("samples/use/hegood.txt")
	if err != nil {
		return "ERROR", errors.New("error in choosing a paragraph")
	}
	return string(b), nil
}
