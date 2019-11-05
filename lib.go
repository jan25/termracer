package main

import (
	"errors"
	"io/ioutil"
	"math/rand"

	"github.com/jan25/termracer/db"
	"github.com/jan25/termracer/pkg/wordwrap"
)

// ChooseParagraph calls the server to fetch a paragraph
func ChooseParagraph() (string, error) {
	samplesFname, err := GetSamplesFilePath()
	if err != nil {
		return fallbackToDefaultParagraph(err)
	}
	samples, err := db.GetSamplesJSON(samplesFname)
	if err != nil {
		return fallbackToDefaultParagraph(err)
	}

	ri := rand.Int() % len(samples)
	p := samples[ri]
	return p.Content, nil
}

// TODO clean this after we're sure fallback isn't necessary
func fallbackToDefaultParagraph(err error) (string, error) {
	Logger.Info("Error reading paragraph from local FS " + err.Error())
	Logger.Info("Falling back to default paragraph")
	// TODO do we need this?
	// Make sure this won't happen
	return firstParagraph, nil
}

// GenerateLocalParagraphs checks if samples/use has > 0 paragraphs
// available. If not tries to generate them
func GenerateLocalParagraphs() error {
	d, _ := GetSamplesUseDir()
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return errors.New("failed to read use directory")
	}
	if len(files) == 0 {
		// TODO generate paragraphs if none available
	}
	return nil
}

// AddNewLines adds new line char to certian words
// to wrap and align the words into seperate lines
func AddNewLines(words []string, width int) {
	processed := []wordwrap.Word{}
	for _, w := range words {
		processed = append(processed, wordwrap.Word{
			Len: len(w),
		})
	}

	wordwrap.Wrap(processed, width)
	for i, w := range processed {
		if w.Wrap {
			words[i] = words[i] + "\n"
		}
	}
}

const firstParagraph = `
She sank more and more into uneasy delirium. At times she shuddered,
turned her eyes from side to side, recognised everyone for a minute,
but at once sank into delirium again. Her breathing was hoarse and
difficult, there was a sort of rattle in her throat.
`
