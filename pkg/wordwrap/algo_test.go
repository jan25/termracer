package wordwrap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	// Test taken from
	// https://en.wikipedia.org/wiki/Line_wrap_and_word_wrap#Minimum_raggedness
	words := strings.Fields("AAA BB CC DDDDD")
	inp := []Word{}
	for _, w := range words {
		inp = append(inp, Word{
			Len:  len(w),
			Wrap: false,
		})
	}

	Wrap(inp, 6)
	// Should be split like so
	// ------    Line width: 6
	// AAA       Remaining space: 3
	// BB CC     Remaining space: 1
	// DDDDD     Remaining space: 1
	assert.True(t, inp[0].Wrap)
	assert.True(t, inp[2].Wrap)
}
