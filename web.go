package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type Website struct {
	LinkedSites []*Website
	Loading bool
	Retrieved bool
	Url string
	Protocol string
}

var urlRegex = regexp.MustCompile(`^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

func (w *Website) RetrieveLinkedSites() {
	body, err := w.get()
	if err == nil {

	}
}

func (w *Website) get() (string, error) {
	response, err := http.Get(w.Url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (w *Website) findOutgoingLinks(body [][]string) {
	return urlRegex.FindAllStringSubmatch(body, 4)
}
