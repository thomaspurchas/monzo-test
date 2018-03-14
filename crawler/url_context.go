package crawler

import (
	"fmt"
	"net/url"
)

type URLContext struct {
	URL           *url.URL
	NormalisedURL *url.URL

	SourceURL           *url.URL
	NormalisedSourceURL *url.URL
}

func (u *URLContext) String() string {
	return fmt.Sprintf("{URL: %v, nURL: %v, sURL: %v, snURL: %v}", u.URL, u.NormalisedURL, u.SourceURL, u.NormalisedSourceURL)
}
