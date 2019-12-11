package db

import (
	"log"
	"math/rand"
	"os"

	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/pkg/wordwrap"
)

// ChooseParagraph calls the server to fetch a paragraph
func ChooseParagraph() string {
	samplesFname, err := config.GetSamplesFilePath()
	if err != nil {
		log.Fatal(err)
	}
	samples, err := GetSamplesJSON(samplesFname)
	if err != nil {
		log.Fatal(err)
	}

	// FIXME rand.Seed(time.Time.Now().Unix())
	ri := rand.Int() % len(samples)
	p := samples[ri]
	return p.Content
}

// GenerateLocalParagraphs checks if samples/use has paragraphs to race with
// available. If not available this func tries to generate them
func GenerateLocalParagraphs() error {
	// TODO Add Loading... thingy before opening UI
	samplesFname, err := config.GetSamplesFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(samplesFname)
	if err != nil && os.IsNotExist(err) {
		err = DownloadSamplesToLocalFS(samplesFname)
		return err
	}

	// We already have the file generated
	return nil
}

// AddNewLines adds new line char to certian words
// to wrap and align the words into seperate lines
func AddNewLines(words []string, width int) int {
	processed := []wordwrap.Word{}
	for _, w := range words {
		processed = append(processed, wordwrap.Word{
			Len: len(w),
		})
	}

	wordwrap.Wrap(processed, width)
	lines := 1
	for i, w := range processed {
		if w.Wrap {
			words[i] = words[i] + "\n"
			lines++
		}
	}
	return lines
}

const firstParagraph = `
She sank more and more into uneasy delirium. At times she shuddered,
turned her eyes from side to side, recognised everyone for a minute,
but at once sank into delirium again. Her breathing was hoarse and
difficult, there was a sort of rattle in her throat.
`
