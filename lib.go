package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jan25/termracer/pkg/wordwrap"
	"github.com/jan25/termracer/server/client"
)

// ChooseParagraph calls the server to fetch a paragraph
func ChooseParagraph() (string, error) {
	p, err := client.ChooseParagraph()
	if err != nil {
		Logger.Info("Error fetching paragraph from server " + err.Error())
		Logger.Info("Falling back to default paragraph")
		// TODO remove this temporary fix
		// returning const paragraph if couldn't reach server
		return firstParagraph, nil
	}
	return p, nil
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
		Logger.Info(fmt.Sprint(words[i]))
	}
}

const firstParagraph = `
She sank more and more into uneasy delirium. At times she shuddered,
turned her eyes from side to side, recognised everyone for a minute,
but at once sank into delirium again. Her breathing was hoarse and
difficult, there was a sort of rattle in her throat.
`
