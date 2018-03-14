package crawler

import (
	"context"
	"log"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/purell"
)

type Page struct {
	url      string
	contents string
}

type vistedURLs struct {
	urls map[string]struct{}
	mux  sync.Mutex
}

type Crawler struct {
	workers         map[string]*worker
	visitedHosts    map[string]struct{}
	externalResults chan *URLContext
	wg              sync.WaitGroup

	Results chan *URLContext

	ctx context.Context
	opt *Options
}

func NewCrawler(ctx context.Context, opt *Options) *Crawler {
	c := &Crawler{}
	c.workers = make(map[string]*worker)
	c.visitedHosts = make(map[string]struct{})

	c.Results = make(chan *URLContext)
	c.externalResults = make(chan *URLContext)

	c.ctx = ctx
	c.opt = opt
	return c
}

func (c *Crawler) Crawl(seed string) {
	defer close(c.Results)
	c.startSeedWorker(seed)

	done := c.workersDone()
	for {
		select {
		case u := <-c.externalResults:
			if c.opt.MultipleDomains {
				if _, exists := c.workers[u.NormalisedURL.Hostname()]; !exists {
					c.startWorker(u)
				}
			}
		case <-done:
			return
		}
	}
}

func (c *Crawler) startWorker(uctx *URLContext) {
	w := newWorker(c.ctx, c.opt, uctx)

	log.Printf("Starting worker for: %s\n", w.host)
	c.wg.Add(1)
	go func() {
		w.run()
		c.wg.Done()
	}()

	c.workers[w.host] = w
	c.addResultsChannel(w.results)
	c.addExternalChannel(w.externalResults)
}

func (c *Crawler) startSeedWorker(seed string) {
	u, _ := url.Parse(seed)
	uctx := &URLContext{URL: u,
		NormalisedURL: cleanURL(u, c.opt.NormalisationFilters)}

	c.startWorker(uctx)
}

func (c *Crawler) workersDone() <-chan struct{} {
	out := make(chan struct{})
	go func() {
		c.wg.Wait()
		out <- struct{}{}
	}()
	return out
}

func (c *Crawler) addExternalChannel(in <-chan *URLContext) {
	go func() {
		quit := c.ctx.Done()
		for n := range in {
			select {
			case c.externalResults <- n:
			case <-quit:
				return
			}

		}
	}()
}

func (c *Crawler) addResultsChannel(in <-chan *URLContext) {
	go func() {
		quit := c.ctx.Done()
		for n := range in {
			select {
			case c.Results <- n:
			case <-quit:
				return
			}
		}
	}()
}

func cleanURL(u *url.URL, flags purell.NormalizationFlags) *url.URL {
	s := purell.NormalizeURL(u, flags)
	u, _ = url.Parse(s)
	return u
}
