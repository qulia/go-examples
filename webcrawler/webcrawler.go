package webcrawler

import (
	"sync"

	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/concurrency/unique"
	"github.com/qulia/go-qulia/lib/set"
	log "github.com/sirupsen/logrus"
)

type UrlData struct {
	Url       string
	SourceUrl string
}

type WebCrawler struct {
	numberOfWorkers int
	urlCount        int
	siteParser      siteparser.Interface
	jobsQueue       chan *UrlData
	urlsAccessor    *unique.Unique[set.Set[string]]
}

func NewWebCrawler(numWorkers, urlCount int, siteParser siteparser.Interface) *WebCrawler {
	wc := WebCrawler{
		numberOfWorkers: numWorkers,
		urlCount:        urlCount,
		siteParser:      siteParser,
		jobsQueue:       make(chan *UrlData, 10),
		urlsAccessor:    unique.NewUnique(set.NewSet[string]()),
	}
	return &wc
}

// Visits BFS manner starting at the startUrl and posts the urls
func (wc *WebCrawler) Visit(startUrl string) set.Set[string] {
	urls, ok := wc.urlsAccessor.Acquire()
	if !ok {
		return urls
	}

	startUrlData := UrlData{
		Url:       startUrl,
		SourceUrl: "",
	}
	urls.Add(startUrlData.Url)
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
	return urls
}

// Pick a url from the job queue, starts a worker to parse, puts this url to output
func (wc *WebCrawler) urlWorker(id int, urlCounter *sync.WaitGroup) {
	log.Infof("At worker %d", id)
	defer log.Infof("Exiting worker %d", id)
	for {
		curUrlData := <-wc.jobsQueue
		log.Infof("Received at worker id:%d urlData:%s", id, curUrlData)
		go wc.parserWorker(curUrlData, urlCounter)
	}
}

// Picks a url to parse, puts new urls in the job queue, if not processed already
func (wc *WebCrawler) parserWorker(urlData *UrlData, urlCounter *sync.WaitGroup) {
	newUrls := wc.siteParser.Parse(urlData.Url)
	log.Infof("Parsed url:%s result:%s", urlData, newUrls)
	urls, ok := wc.urlsAccessor.Acquire()
	if !ok {
		return
	}
	defer wc.urlsAccessor.Release()
	for _, newUrl := range newUrls {
		if !urls.Contains(newUrl) {
			urls.Add(newUrl)
			urlCounter.Done()
			wc.jobsQueue <- &UrlData{
				Url:       newUrl,
				SourceUrl: urlData.Url,
			}
		}
	}
}

func (wc *WebCrawler) Stop() {
	wc.urlsAccessor.Close()
}
