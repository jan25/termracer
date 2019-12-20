package data

import (
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jan25/termracer/config"
	"github.com/jan25/termracer/pkg/utils"
)

// EnsureDataDirs checks to see data dirs required for application are present
// creates dirs/files if not present
func EnsureDataDirs() error {
	// ensure top level directory
	tld, err := config.GetTopLevelDir()
	if err != nil {
		return err
	}
	if err := utils.CreateDirIfNotExists(tld); err != nil {
		return err
	}

	if err := generateLocalParagraphs(); err != nil {
		return err
	}

	// ensure racehistory file
	rh, err := config.GetHistoryFilePath()
	if err != nil {
		return err
	}
	if err := utils.CreateFileIfNotExists(rh); err != nil {
		return err
	}

	return nil
}

// DownloadSamplesToLocalFS downloads and stores samples.json file locally
func downloadSamplesToLocalFS(fname string) error {
	url := "https://github.com/jan25/termracer/raw/master/data/scripts/samples.gz"
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
