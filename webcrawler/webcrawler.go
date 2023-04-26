package webcrawler

import (
	"net/url"
	"sync"
	"time"

	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/algo/ratelimiter"
	"github.com/qulia/go-qulia/algo/ratelimiter/leakybucket"
	"github.com/qulia/go-qulia/concurrency/unique"
	"github.com/qulia/go-qulia/lib/set"
	log "github.com/sirupsen/logrus"
)

type UrlData struct {
	Url       *url.URL
	SourceUrl *url.URL
}

type WebCrawler struct {
	numberOfWorkers           int
	urlCount                  int
	siteParser                siteparser.Interface
	jobsQueue                 chan *UrlData
	urlsAccessor              *unique.Unique[set.Set[url.URL]]
	hostController            map[string]ratelimiter.RateLimiterBuffered
	perHostRequestGapDuration time.Duration
}

func NewWebCrawler(numWorkers, urlCount int, siteParser siteparser.Interface) *WebCrawler {
	wc := WebCrawler{
		numberOfWorkers:           numWorkers,
		perHostRequestGapDuration: time.Second * 2,
		urlCount:                  urlCount,
		siteParser:                siteParser,
		jobsQueue:                 make(chan *UrlData, 10),
		urlsAccessor:              unique.NewUnique(set.NewSet[url.URL]()),
		hostController:            make(map[string]ratelimiter.RateLimiterBuffered),
	}
	return &wc
}

// Visits BFS manner starting at the startUrl and posts the urls
func (wc *WebCrawler) Visit(startUrl string) (set.Set[url.URL], error) {
	urls, ok := wc.urlsAccessor.Acquire()
	if !ok {
		return urls, nil
	}

	urlParsed, err := url.Parse(startUrl)
	if err != nil {
		return *new(set.Set[url.URL]), err
	}
	startUrlData := UrlData{
		Url:       urlParsed,
		SourceUrl: nil,
	}
	wc.jobsQueue <- &startUrlData
	wc.urlsAccessor.Release()

	urlCounter := &sync.WaitGroup{}
	urlCounter.Add(wc.urlCount)
	// Start workers
	log.Infof("Starting %d workers...", wc.numberOfWorkers)
	for i := 0; i < wc.numberOfWorkers; i++ {
		go wc.urlWorker(i, urlCounter)
	}

	urlCounter.Wait()
	urls, _ = wc.urlsAccessor.Acquire()
	wc.urlsAccessor.Release()
	return urls, err
}

// Pick a url from the job queue, starts a worker to parse, puts this url to output
func (wc *WebCrawler) urlWorker(id int, urlCounter *sync.WaitGroup) {
	log.Infof("At worker %d", id)
	defer log.Infof("Exiting worker %d", id)
	for {
		curUrlData := <-wc.jobsQueue
		log.Infof("Received at worker id:%d urlData:%s", id, curUrlData)
		rch, allowed, ok := wc.allow(*curUrlData.Url)
		if !ok {
			return
		}
		if !allowed {
			go func(urlData *UrlData) {
				time.Sleep(wc.perHostRequestGapDuration)
				wc.jobsQueue <- curUrlData
			}(curUrlData)
		} else {
			go func(rch <-chan interface{}) {
				<-rch
				wc.parserWorker(curUrlData, urlCounter)
			}(rch)
		}
	}
}

// Picks a url to parse, puts new urls in the job queue, if not processed already
func (wc *WebCrawler) parserWorker(urlData *UrlData, urlCounter *sync.WaitGroup) {
	urls, ok := wc.urlsAccessor.Acquire()
	if !ok {
		return
	}

	if urls.Len() == wc.urlCount || urls.Contains(*urlData.Url) {
		wc.urlsAccessor.Release()
		return
	}

	urls.Add(*urlData.Url)
	urlCounter.Done()
	wc.urlsAccessor.Release()

	newUrls := wc.siteParser.Parse(urlData.Url.String())
	log.Infof("Parsed url:%s result:%s", urlData, newUrls)

	urls, ok = wc.urlsAccessor.Acquire()
	if !ok {
		return
	}
	defer wc.urlsAccessor.Release()
	for _, newUrl := range newUrls {
		urlParsed, _ := url.Parse(newUrl)
		if !urls.Contains(*urlParsed) && urls.Len() < wc.urlCount {
			wc.jobsQueue <- &UrlData{
				Url:       urlParsed,
				SourceUrl: urlData.Url,
			}
		}
	}
}

func (wc *WebCrawler) allow(cur url.URL) (<-chan interface{}, bool, bool) {
	_, ok := wc.urlsAccessor.Acquire()
	if !ok {
		return nil, false, false
	}
	defer wc.urlsAccessor.Release()
	if wc.hostController[cur.Host] == nil {
		wc.hostController[cur.Host] = leakybucket.NewLeakyBucket(1, 1, wc.perHostRequestGapDuration)
	}
	rl := wc.hostController[cur.Host]
	rch, ok := rl.Allow()
	if !ok {
		log.Infof("host is rate limited: %s", cur.Host)
		return nil, false, true
	}
	return rch, true, true
}

func (wc *WebCrawler) Stop() {
	wc.urlsAccessor.Close()
}
