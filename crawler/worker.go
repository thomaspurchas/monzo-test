package crawler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/temoto/robotstxt"
)

const (
	robotsPath    = "/robots.txt"
	defaultRobots = "User-agent: *\nDisallow:"
)

type worker struct {
	host string
	seed *URLContext

	fetchResults    chan *fetchResult
	visitedURLs     map[string]struct{}
	waitGroup       sync.WaitGroup
	results         chan *URLContext
	externalResults chan *URLContext

	robots *robotstxt.Group

	ctx context.Context

	opt *Options
}

func newWorker(ctx context.Context, opt *Options, seed *URLContext) *worker {
	r := &worker{}
	r.seed = seed
	r.host = seed.NormalisedURL.Hostname()

	r.fetchResults = make(chan *fetchResult)
	r.visitedURLs = make(map[string]struct{})
	r.results = make(chan *URLContext)
	r.externalResults = make(chan *URLContext)

	r.ctx = ctx
	r.opt = opt

	return r
}

func (w *worker) run() {
	defer close(w.results)
	defer close(w.externalResults)
	defer log.Printf("Worker for '%s' finished \n", w.host)

	w.waitGroup.Add(1)
	go func() {
		defer w.waitGroup.Done()

		f := &fetchResult{}
		f.url = w.seed
		f.foundURLs = []*url.URL{w.seed.NormalisedURL}
		f.resp = make(chan struct{})

		select {
		case w.fetchResults <- f:
		case <-w.ctx.Done():
		}
		select {
		case <-f.resp:
		case <-w.ctx.Done():
		}
	}()

	fetchDone := w.fetchDone()

	for {
		quit := w.ctx.Done()
		select {
		case <-quit:
			// Exit if we are told to
			log.Printf("Worker for '%s' stopping due to quit signal\n", w.host)
			return
		case r := <-w.fetchResults:
			// Return successfully fetched page to the crawler.
			// Only returning the URL because this challenge doesn't need the page
			// contents
			select {
			case w.results <- r.url:
			case <-quit:
				return
			}

			internalu, externalu := w.processResults(r)
			for _, u := range internalu {
				w.fetchURL(u)
			}
			for _, u := range externalu {
				// Return found urls that are external to this host, the crawler
				// is then responsible for doing something with them.
				select {
				case w.externalResults <- u:
				case <-quit:
					return
				}
			}
			// Send resp to fetcher, allowing it to die
			r.resp <- struct{}{}
		case <-fetchDone:
			log.Printf("Finished fetching from host: %s\n", w.host)
			return
		}
	}
}

func (w *worker) fetchDone() chan struct{} {
	c := make(chan struct{})
	go func() {
		w.waitGroup.Wait()
		log.Println("All fetchers done")
		c <- struct{}{}
	}()
	return c
}

func (w *worker) fetchURL(u *URLContext) {
	if w.robotsAllowed(u) {
		log.Printf("Fetching: %s\n", u)

		w.waitGroup.Add(1)
		go fetchURL(w.ctx, w, u)
	} else {
		log.Printf("Disallowed by robots: %s\n", u.NormalisedURL)
	}
}

func (w *worker) robotsAllowed(u *URLContext) bool {
	if w.robots == nil {
		log.Println("Fetching robots.txt")
		var robots *robotstxt.RobotsData
		robotsURL, _ := url.Parse(robotsPath)

		res, err := http.Get(u.NormalisedURL.ResolveReference(robotsURL).String())
		if err == nil {
			robots, err = robotstxt.FromResponse(res)
			if err != nil {
				robots = nil
			}
		}

		if robots == nil {
			robots, _ = robotstxt.FromString(defaultRobots)
		}

		w.robots = robots.FindGroup(w.opt.AgentName)
	}

	return w.robots.Test(u.NormalisedURL.EscapedPath())
}

func (w *worker) processResults(r *fetchResult) ([]*URLContext, []*URLContext) {
	// Filters urls for visited urls
	urls := make([]*URLContext, 0, len(r.foundURLs))
	var externalUrls []*URLContext
	for _, u := range r.foundURLs {
		urlctx := &URLContext{}
		urlctx.URL = u
		urlctx.NormalisedURL = w.cleanURL(u)
		urlctx.SourceURL = r.url.URL
		urlctx.NormalisedSourceURL = r.url.NormalisedURL

		if urlctx.NormalisedURL.Hostname() != w.host {
			externalUrls = append(externalUrls, urlctx)
		} else {
			if _, exists := w.visitedURLs[urlctx.NormalisedURL.String()]; !exists {
				urls = append(urls, urlctx)
				// Mark the page as visited, because it will be visited at some
				// point in the future
				w.visitedURLs[urlctx.NormalisedURL.String()] = struct{}{}
			}
		}
	}

	return urls, externalUrls
}

func (w *worker) cleanURL(u *url.URL) *url.URL {
	return cleanURL(u, w.opt.NormalisationFilters)
}
