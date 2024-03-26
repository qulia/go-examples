package webcrawler_test

import (
	"testing"
	"time"

	"github.com/qulia/go-examples/webcrawler"
	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/lib/set"
	"github.com/stretchr/testify/assert"
)

func TestWebCrawlerBasic(t *testing.T) {
	urlMap := map[string][]string{
		"http://site1.com/page1": {"http://site3.com/page1", "http://site4.com/page1"},
		"http://site2.com/page1": {"http://site2.com/page1", "http://site4.com/page1"},
		"http://site3.com/page1": {"http://site2.com/page1", "http://site6.com/page1"},
		"http://site4.com/page1": {
			"http://site1.com/page1",
			"http://site2.com/page1",
			"http://site3.com/page1",
			"http://site4.com/page1",
			"http://site6.com/page1",
		},

		"http://site6.com/page1": {
			"http://site1.com/page2",
			"http://site1.com/page3",
			"http://site1.com/page4",
			"http://site1.com/page5",
		},
	}

	expectedUrls := set.NewSet[string]().FromSlice([]string{
		"http://site1.com/page1",
		"http://site2.com/page1",
		"http://site3.com/page1",
		"http://site4.com/page1",
		"http://site6.com/page1",
		"http://site1.com/page2",
		"http://site1.com/page3",
		"http://site1.com/page4",
		"http://site1.com/page5",
	})

	sp := siteparser.NewMockSiteParser(urlMap, nil)

	wc := webcrawler.NewWebCrawler(4, expectedUrls.Len(), time.Minute, time.Second*2, sp)
	foundUrls, err := wc.Visit("http://site1.com/page1")
	assert.Equal(t, err.Error(), "")
	t.Logf("Found urls %v", foundUrls)
	foundUrlsSet := set.NewSet[string]()
	for _, u := range foundUrls.ToSlice() {
		foundUrlsSet.Add(u)
	}
	assert.True(t, expectedUrls.IsSubsetOf(foundUrlsSet) && expectedUrls.IsSupersetOf(foundUrlsSet))

	wc.Stop()
}
