package main

import (
	"fmt"

	"github.com/thomaspurchas/monzo-test/crawler"
)

func main() {
	pages := crawler.Crawl("http://monzo.com")
	for page := range pages {
		fmt.Println(page)
	}
}
