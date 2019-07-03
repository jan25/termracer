package client

import (
	"io/ioutil"
	"net/http"
)

const (
	defaultHost  = "http://localhost:8080"
	paraEndpoint = "/paragraph"
)

// ChooseParagraph is a api wrapper to fetch
// paragraph from server
func ChooseParagraph() (string, error) {
	uri := defaultHost + paraEndpoint
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(p), nil
}
