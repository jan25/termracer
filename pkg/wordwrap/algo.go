package wordwrap

import (
	"errors"
)

// Word represents a string of consecutive chars in a line
type Word struct {
	// Len represents length of the word
	Len int
	// Wrap is true if this word is last one in current line
	Wrap bool
}

// Wrap executes the algorithm to wrap list of words
// algorithm: https://en.wikipedia.org/wiki/Line_wrap_and_word_wrap#Minimum_raggedness
func Wrap(words []Word, width int) error {
	lens := []int{}
	for _, w := range words {
		if w.Len > width {
			return errors.New("One of words is longer than line width")
		}
		// reset word wraps
		w.Wrap = false

		lens = append(lens, w.Len)
	}
	for i := len(lens) - 2; i >= 0; i-- {
		lens[i] += lens[i+1]
	}

	// Call dynamic programming algorithm and
	// mark last words in each line
	nexti := map[int]int{}
	dp(0, lens, width, map[int]int{}, nexti)
	i, ok := 0, true
	for {
		i, ok = nexti[i]
		if !ok {
			break
		}
		words[i].Wrap = true
		i++
	}

	return nil
}

// lens: suffixed cumulative sum of word lengths
// next: links indices of first words in each line
func dp(i int, lens []int, w int, memo map[int]int, next map[int]int) int {
	n := len(lens)
	if i >= n || lens[i]+(n-i)-1 <= w {
		return 0
	}
	if v, ok := memo[i]; ok {
		return v
	}
	k, sol := i, w*w*n
	for j := i + 1; j < n; j++ {
		chars := (lens[i] - lens[j]) + (j - i)
		if chars > w {
			break
		}
		k = j
		spac := w - chars
		subsol := spac*spac + dp(j, lens, w, memo, next)
		if subsol < sol {
			sol = subsol
			k = j
		}
	}
	next[i] = k - 1
	memo[i] = sol
	return sol
}
