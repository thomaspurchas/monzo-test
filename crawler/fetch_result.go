package crawler

import (
	"fmt"
	"net/url"
)

type fetchResult struct {
	url       *URLContext
	foundURLs []*url.URL
	resp      chan struct{}
}

func (f *fetchResult) String() string {
	return fmt.Sprintf("{url: %v, foundURLs: %v}", f.url, f.foundURLs)
}
