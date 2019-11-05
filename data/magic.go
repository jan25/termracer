// This is one-off script. Used mainly to generate samples.gz in this direcgtory
// Fetches https://github.com/jan25/wpm/tree/master/wpm/data/samples.json.gz
// and parses it into custom JSON format, creates samples.gz
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jan25/termracer/db"
)

func main() {
	fmt.Println("Yes this is the one!")
	err := getAndSaveData()
	if err != nil {
		fmt.Println("Sorry, something went wrong. Here it is: \n", err)
	} else {
		fmt.Println("Done. samples.gz saved :)")
	}
}

func getAndSaveData() error {
	url := "https://github.com/jan25/wpm/blob/master/wpm/data/examples.json.gz?raw=true"

	client := new(http.Client)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept-Encoding", "gzip")
	resp, err := client.Do(request)
	if err != nil {
		return errors.New("Failed to GET remote file. " + err.Error())
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return errors.New("Failed to create gzip.NewReader. " + err.Error())
	}
	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("Failed to read response to bytes. " + err.Error())
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

	err = ioutil.WriteFile(
		"./samples.gz",
		bytes,
		0644,
	)
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

// based on https://gist.github.com/alex-ant/aeaaf497055590dacba760af24839b8d
func unzip(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return resB.Bytes(), nil
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
