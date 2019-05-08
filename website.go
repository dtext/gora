package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Website struct {
	LinkedSites map[*Website]int
	Loading     bool
	Retrieved   bool
	Url         string
	lock        sync.Mutex
	err         error
}

var urlRegex = regexp.MustCompile(`https?://([a-zA-Z0-9]*\.)([a-zA-Z0-9]*\.?)*/?`)

func (w *Website) acquireOkForLoad() bool {
	// acquire lock, release when done
	w.lock.Lock()
	defer w.lock.Unlock()

	// website is ok for load if it's neither retrieved nor loading
	ok := !w.Retrieved && !w.Loading
	if !ok {
		return false
	}
	w.Loading = true
	return true
}

func (w *Website) doneLoading() {
	// acquire lock, release when done
	w.lock.Lock()
	defer w.lock.Unlock()

	// set variables
	w.Loading = false
	w.Retrieved = true

}

func (w *Website) RetrieveAsync(theWeb *Web, results chan *Website, work chan *Website) {
	// get HTML, filter out URLS, build Website structs
	w.RetrieveLinks(theWeb)

	// done fetching: send to results channel
	results <- w

	// send unfetched links to work channel
	for website := range w.LinkedSites {
		if !website.acquireOkForLoad() {
			continue
		}
		work <- website
	}
}

func (w *Website) RetrieveLinks(theWeb *Web) {
	// HTTP(S) GET
	body, err := w.get()
	if err != nil {
		w.err = err
		return
	}

	// find outgoing URLs
	outgoingLinks := urlRegex.FindAllString(body, -1)
	for _, link := range outgoingLinks {
		// do we have that URL already?
		website, ok := (*theWeb)[link]

		// if not, create a Website struct for it and save that
		if !ok {
			website = &Website{
				LinkedSites: make(map[*Website]int),
				Loading:     false,
				Retrieved:   false,
				Url:         link,
				lock:        sync.Mutex{},
				err:         nil,
			}
			(*theWeb)[link] = website
		}

		// add to linked sites, counting links
		linkCount, exists := w.LinkedSites[website]
		if !exists {
			// initialize counter
			w.LinkedSites[website] = 0
			linkCount = 0
		}
		w.LinkedSites[website] = linkCount + 1
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

func (w Website) String() string {
	var sb strings.Builder
	sb.WriteString(w.Url)
	sb.WriteString(" links the following websites:\n")
	for link := range w.LinkedSites {
		sb.WriteString("    ")
		sb.WriteString(link.Url)
		sb.WriteString("\n")
	}
	return sb.String()
}
