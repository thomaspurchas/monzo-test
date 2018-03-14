package main

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/PuerkitoBio/purell"
	"github.com/gonum/graph/encoding/dot"

	"github.com/thomaspurchas/monzo-test/crawler"
	"github.com/thomaspurchas/monzo-test/grapher"
)

const dotHeader = `rankdir="LR";
concentrate=true;
splines=polyline;
ranksep="2 equally"\n`

func pass(u *url.URL) *url.URL {
	return u
}

func main() {
	ctx := context.Background()

	opt := &crawler.Options{
		AgentName:            "test",
		URLFilters:           []crawler.URLFilter{pass},
		NormalisationFilters: purell.FlagsUsuallySafeGreedy | purell.FlagRemoveFragment}

	c := crawler.NewCrawler(ctx, opt)
	go c.Crawl("http://monzo.com")

	g := grapher.BuildGraph(c.Results)
	data, err := dot.Marshal(g, "Monzo", "", "  ", false)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("out.dot", data, 0644)
}
