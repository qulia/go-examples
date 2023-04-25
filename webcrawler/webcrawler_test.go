package webcrawler_test

import (
	"testing"

	"github.com/qulia/go-examples/webcrawler"
	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/lib/set"
	"github.com/stretchr/testify/assert"
)

func TestWebCrawlerBasic(t *testing.T) {
	urlMap := map[string][]string{
		"site1": {"site3", "site4"},
		"site2": {"site2", "site4"},
		"site3": {"site2", "site6"},
		"site4": {"site1", "site2", "site3", "site4", "site6"},
	}
	sp := siteparser.NewMockSiteParser(urlMap)

	wc := webcrawler.NewWebCrawler(4, 4, sp)
	foundUrls := wc.Visit("site1")
	t.Logf("Found urls %v", foundUrls)

	expectedUrls := set.NewSet[string]().FromSlice([]string{"site1", "site2", "site3", "site4", "site6"})
	assert.True(t, expectedUrls.IsSubsetOf(foundUrls) && expectedUrls.IsSupersetOf(foundUrls))
}
