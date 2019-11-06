package db

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jan25/termracer/pkg/utils"
)

// Sample is a paragraph sample used for races
type Sample struct {
	Content string `json:"content"`
}

// GetSamplesJSON returns the JSON file contents
func GetSamplesJSON(fname string) ([]Sample, error) {
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

// DownloadSamplesToLocalFS downloads stores samples.json file locally
func DownloadSamplesToLocalFS(fname string) error {
	// FIXME add URL that would be good after merging MR
	url := "https://github.com/jan25/termracer/blob/f5fad2102084b2f4cd8f9dcb618af6df2d9b9d84/data/samples.gz?raw=true"
	bytes, err := DownloadGzipFile(url)
	if err != nil {
		return err
	}
	if err = utils.CreateFileIfNotExists(fname); err != nil {
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
