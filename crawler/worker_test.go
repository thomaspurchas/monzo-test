package crawler

import (
	"net/url"
	"testing"

	"github.com/PuerkitoBio/purell"
)

const pFlags = purell.FlagsUsuallySafeGreedy | purell.FlagRemoveFragment

func TestProcessResults(t *testing.T) {
	testURL, _ := url.Parse("http://thingy.com")
	u := &URLContext{
		URL:           testURL,
		NormalisedURL: cleanURL(testURL, pFlags)}

	r := &fetchResult{url: u,
		foundURLs: []*url.URL{u.NormalisedURL}}

	w := &worker{host: testURL.Hostname()}
	w.opt = &Options{NormalisationFilters: pFlags}

	in, out := w.processResults(r)

	t.Log(r.foundURLs)
	t.Log(in)

	if in[0].URL.String() != testURL.String() {
		t.Errorf("Internal url returned does not match testURL: %v != %v", in[0].URL, testURL)
	}
	if len(out) > 0 {
		t.Error("External url count is not 0!")
	}

	externalURL, _ := url.Parse("http://external.com")
	r.foundURLs = append(r.foundURLs, externalURL)

	in, out = w.processResults(r)

	t.Log(in)
	t.Log(out)

	if len(out) != 1 || out[0].URL.String() != externalURL.String() {
		t.Errorf("External host not filtered!")
	}
}
