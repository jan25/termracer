package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func chooseParagraph() (string, error) {
	files, err := ioutil.ReadDir("samples/use")
	n := len(files)
	if err != nil || n == 0 {
		return "ERROR", errors.New("failed to read use directory")
	}

	rand.Seed(time.Now().Unix())
	randf := files[rand.Int31n(int32(n))]

	b, err := ioutil.ReadFile("samples/use/" + randf.Name())
	if err != nil {
		return "ERROR", errors.New("error in reading a paragraph file")
	}
	return string(b), nil
}

func paragraphHandler(w http.ResponseWriter, r *http.Request) {
	p, err := chooseParagraph()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, p)
}

func main() {
	http.HandleFunc("/paragraph", paragraphHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
