package utils

import (
	"errors"
	"math"
)

// CalculateWpm calculates words per minute
// based on words typed so far and time elapsed
func CalculateWpm(countCh int, secs float64) int {
	// average word length is 5
	return int((60 * float64(countCh) / 5) / secs)
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
