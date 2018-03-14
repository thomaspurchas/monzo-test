package crawler

import (
	"net/url"

	"github.com/PuerkitoBio/purell"
)

type URLFilter func(*url.URL) *url.URL

type Options struct {
	URLFilters           []URLFilter
	NormalisationFilters purell.NormalizationFlags
	MultipleDomains      bool

	AgentName string
}
