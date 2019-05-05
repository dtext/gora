package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Website struct {
	LinkedSites []*Website
	Loading     bool
	Retrieved   bool
	Url         string
	lock        sync.Mutex
	err         error
}

type Web map[string]Website

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

func BuildGraph(startingUrl string, limit, timeout int, resultHandler func(*Website)) Web {
	theWeb := make(Web)
	results := make(chan *Website)
	work := make(chan *Website)

	w := Website{LinkedSites: []*Website{},
		Loading:   true,
		Retrieved: false,
		Url:       startingUrl,
		err:       nil,
		lock:      sync.Mutex{},
	}
	theWeb[startingUrl] = w

	// Timeout after <timeout> seconds or never
	var timeoutChannel = make(chan time.Time)
	go func() {
		if timeout == 0 {
			return
		}
		time.Sleep(time.Duration(timeout) * time.Second)
		timeoutChannel <- time.Now()
	}()

	// Kickoff
	go w.RetrieveAsync(&theWeb, results, work)

	// Wait for timeout, more websites to fetch and results
	select {
	case <-timeoutChannel:
		// TODO stop goroutines
		break
	case website := <-work:
		if limit == 0 {
			break
		}
		limit--
		go website.RetrieveAsync(&theWeb, results, work)
	case res := <-results:
		resultHandler(res)
	}
	return theWeb
}

func (w *Website) RetrieveAsync(theWeb *Web, results chan *Website, work chan *Website) {
	// get HTML, filter out URLS, build Website structs
	w.RetrieveLinks(theWeb)

	// done fetching: send to results channel
	results <- w

	// send unfetched links to work channel
	for _, website := range w.LinkedSites {
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
			website = Website{
				LinkedSites: make([]*Website, 4),
				Loading:     false,
				Retrieved:   false,
				Url:         link,
				lock:        sync.Mutex{},
				err:         nil,
			}
			(*theWeb)[link] = website
		}
		// TODO use map for multiple link support
		w.LinkedSites = append(w.LinkedSites, &website)
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
