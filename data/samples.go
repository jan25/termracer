package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/jan25/termracer/config"
)

// Sample is a paragraph sample used for races
type Sample struct {
	Content string `json:"content"`
}

// ChooseParagraph calls the server to fetch a paragraph
func ChooseParagraph() string {
	samplesFname, err := config.GetSamplesFilePath()
	if err != nil {
		log.Fatal(err)
	}
	samples, err := getSamplesJSON(samplesFname)
	if err != nil {
		log.Fatal(err)
	}

	// FIXME rand.Seed(time.Time.Now().Unix())
	ri := rand.Int() % len(samples)
	p := samples[ri]
	return p.Content
}

// getSamplesJSON returns the JSON file contents
func getSamplesJSON(fname string) ([]Sample, error) {
	_, err := os.Stat(fname)
	if err != nil && os.IsNotExist(err) {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	var samples []Sample
	if err = json.Unmarshal(bytes, &samples); err != nil {
		return nil, err
	}

	return samples, nil
}

// generateLocalParagraphs checks if samples/use has paragraphs to race with
// available. If not available this func tries to generate them
func generateLocalParagraphs() error {
	// TODO Add Loading... thingy before opening UI
	samplesFname, err := config.GetSamplesFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(samplesFname)
	if err != nil && os.IsNotExist(err) {
		err = downloadSamplesToLocalFS(samplesFname)
		return err
	}

	// We already have the file generated
	return nil
}
