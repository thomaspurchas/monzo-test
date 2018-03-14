package crawler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchURL(ctx context.Context, w *worker, u *URLContext) {
	defer w.waitGroup.Done()
	defer log.Println("Finished processing:", u.URL)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(u.NormalisedURL.String())
	if err != nil {
		return
	}
	log.Printf("Fetched: %s\n", u.URL)

	doc, err := goquery.NewDocumentFromResponse(res)

	if err != nil {
		return
	}

	urls := processDoc(doc)
	fu := make([]*url.URL, 0, len(urls))

	for _, u := range urls {
		if len(w.opt.URLFilters) > 0 {
			for _, filter := range w.opt.URLFilters {
				if f := filter(u); f != nil {
					fu = append(fu, f)
				}
			}
		} else {
			fu = append(fu, u)
		}

	}

	result := &fetchResult{
		url:       u,
		foundURLs: fu,
		resp:      make(chan struct{})}

	select {
	case w.fetchResults <- result:
	case <-w.ctx.Done():
	}
	select {
	case <-result.resp:
	case <-w.ctx.Done():
	}
}

func processDoc(doc *goquery.Document) []*url.URL {
	var base *url.URL
	if baseURL, _ := doc.Find("base[href]").Attr("href"); baseURL != "" {
		base, _ = url.Parse(baseURL)
	} else {
		base = doc.Url
	}

	urls := doc.Find("a[href]").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Attr("href")
		return val
	})

	var result []*url.URL

	for _, u := range urls {
		if ru, err := url.Parse(u); err == nil {
			au := base.ResolveReference(ru)

			result = append(result, au)
		} else {
			log.Printf("Ignored: unable to parse: %s", u)
		}

	}

	return result
}
