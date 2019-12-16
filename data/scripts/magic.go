// This is one-off script. Used mainly to generate samples.gz in this direcgtory
// Fetches https://github.com/jan25/wpm/blob/master/wpm/data/scripts/examples.json.gz?raw=true(I copied this from some other repo, thank them)
// and parses it into custom JSON format, creates samples.gz
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/jan25/termracer/data"
	"github.com/jan25/termracer/pkg/utils"
)

func main() {
	fmt.Println("Starting..")
	err := getAndSaveData()
	if err != nil {
		fmt.Println("Sorry, something went wrong. Here it is: \n", err)
	} else {
		fmt.Println("Done. samples.gz saved :)")
	}
}

func getAndSaveData() error {
	url := "https://github.com/jan25/wpm/blob/master/wpm/data/scripts/examples.json.gz?raw=true"
	bytes, err := db.DownloadGzipFile(url)
	if err != nil {
		return errors.New("Failed to download gzip file. " + err.Error())
	}

	var parsed [][]interface{}
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		return errors.New("Failed to parse unzipped bytes. " + err.Error())
	}

	samples := []db.Sample{}
	for _, s := range parsed {
		sample := db.Sample{
			Content: s[2].(string), // Why 2? refer to sample at end of this script
		}
		samples = append(samples, sample)
	}
	return saveData(samples)
}

func saveData(data []db.Sample) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if bytes, err = zip(bytes); err != nil {
		return nil
	}

	err = utils.WriteToFile("./samples.gz", bytes)
	return err
}

// based on https://gist.github.com/alex-ant/aeaaf497055590dacba760af24839b8d
func zip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err = gz.Flush(); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Sample from github .gz file
// [
// 	 [
// 		"Tetsuro Araki",
// 		"Death Note",
// 		"You're asking me why? I did it 'cause I was bored.",
// 		3550614
// 	 ],
// ]
