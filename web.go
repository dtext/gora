package main

import (
	"sync"
	"time"
)

type Web map[string]*Website

func BuildGraph(startingUrl string, limit, timeout int, resultHandler func(*Website)) Web {
	theWeb := make(Web)
	results := make(chan *Website)
	work := make(chan *Website)

	w := Website{LinkedSites: make(map[*Website]int),
		Loading:   true,
		Retrieved: false,
		Url:       startingUrl,
		err:       nil,
		lock:      sync.Mutex{},
	}
	theWeb[startingUrl] = &w

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
	limit--
	go w.RetrieveAsync(&theWeb, results, work)

	// Wait for timeout, more websites to fetch and results
loop:
	for {
		select {
		case <-timeoutChannel:
			// TODO stop goroutines
			break loop
		case website := <-work:
			if limit == 0 {
				break loop
			}
			limit--
			go website.RetrieveAsync(&theWeb, results, work)
		case res := <-results:
			resultHandler(res)
		}
	}
	return theWeb
}
