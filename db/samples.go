package db

import (
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jan25/termracer/pkg/utils"
)

// Sample is a paragraph sample used for races
type Sample struct {
	Content string `json:"content"`
}

// GetSamplesJSON returns the JSON file contents
func GetSamplesJSON(fname string) ([]Sample, error) {
	// TODO
	return nil, nil
}

// download file and store json file locally
func initSamples(fname string) error {
	url := "" // FIXME add URL
	bytes, err := DownloadGzipFile(url)
	if err != nil {
		return err
	}
	return utils.WriteToFile(fname, bytes)
}

// DownloadGzipFile downloads gzip file from a remote URL
func DownloadGzipFile(url string) ([]byte, error) {
	client := new(http.Client)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept-Encoding", "gzip") // http will auto unzip

	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("Failed to GET remote file. " + err.Error())
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to create gzip.NewReader. " + err.Error())
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("Failed to read response to bytes. " + err.Error())
	}

	return bytes, nil
}
