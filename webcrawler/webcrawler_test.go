package webcrawler_test

import (
	"testing"

	"github.com/qulia/go-examples/webcrawler"
	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/lib"
	"github.com/qulia/go-qulia/utils"
	log "github.com/sirupsen/logrus"
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

	wc := webcrawler.NewWebCrawler(4, sp)

	urls := make(chan webcrawler.UrlData, 1)
	go wc.Visit("site1", urls)

	expectedUrls := []interface{}{"site1", "site2", "site3", "site4", "site6"}
	var foundUrls []interface{}
	// Log urls as they arrive
	for url := range urls {
		if url.Url == "" {
			log.Infof("exiting")
			break
		}
		log.Infof("Url found: %v", url)
		foundUrls = append(foundUrls, url.Url)

		if len(foundUrls) == len(expectedUrls) {
			log.Infof("Found all urls")
			wc.Stop()
			urls <- webcrawler.UrlData{}
		}
	}

	log.Infof("Found urls %v", foundUrls)
	assert.True(t, utils.SliceContains(foundUrls, expectedUrls, lib.HashKeyFunc))
}
