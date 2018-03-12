package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	url      string
	contents string
}

type vistedURLs struct {
	urls map[string]struct{}
	mux  sync.Mutex
}

func filterURLs(urls chan string) chan string {
	c := make(chan string)
	vistedURLs := vistedURLs{urls: make(map[string]struct{})}
	go func() {
		for s := range urls {
			u, err := url.Parse(s)

			if err == nil && u.Hostname() == "monzo.com" {
				if _, ok := vistedURLs.urls[u.String()]; ok == false {
					c <- u.String()
					vistedURLs.urls[u.String()] = struct{}{}
				}
			} else {
				fmt.Printf("Dropping %s", u.String())
			}
		}
	}()
	return c
}

func Crawl(URL string) <-chan Page {
	c := make(chan Page)
	urlsToVisit := make(chan string, 1)

	// vistedURLs := vistedURLs{urls: make(map[string]struct{})}

	urlsToVisit <- URL

	filteredURLs := filterURLs(urlsToVisit)

	for currentURL := range filteredURLs {
		fmt.Println(currentURL)
		res, err := http.Get(currentURL)
		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromResponse(res)

		if err != nil {
			panic(err)
		}

		urls := doc.Find("a[href]").Map(func(_ int, s *goquery.Selection) string {
			val, _ := s.Attr("href")

			u, err := url.Parse(val)

			if err != nil {
				panic(err)
			}

			u = doc.Url.ResolveReference(u)
			return u.String()
		})

		go func(urls []string, urlChan chan string) {
			for _, i := range urls {
				urlChan <- i
			}
		}(urls, urlsToVisit)
	}

	return c
}
