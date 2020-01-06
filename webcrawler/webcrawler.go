package webcrawler

import (
	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/lib"
	"github.com/qulia/go-qulia/lib/queue"
	"github.com/qulia/go-qulia/lib/set"
	log "github.com/sirupsen/logrus"
)

const (
	allUrls = "allUrls"
)

type UrlData struct {
	Url       string
	SourceUrl string
}

type WebCrawler struct {
	numberOfWorkers int
	siteParser      siteparser.Interface
	jobsQueue       chan *queue.Queue
}

func NewWebCrawler(numWorkers int, siteParser siteparser.Interface) *WebCrawler {
	wc := WebCrawler{
		numberOfWorkers: numWorkers,
		siteParser:      siteParser,
		jobsQueue:       make(chan *queue.Queue, 1),
	}

	return &wc
}

// Visits BFS manner starting at the startUrl and posts the urls
func (wc *WebCrawler) Visit(startUrl string, urls chan<- UrlData) {

	q := queue.NewQueue()
	urlSet := set.NewSet(lib.HashKeyFunc)
	startUrlData := UrlData{
		Url:       startUrl,
		SourceUrl: "",
	}
	urlSet.Add(startUrlData.Url)
	q.Metadata[allUrls] = urlSet
	q.Enqueue(startUrlData)

	// Start workers
	log.Infof("Starting %d workers...", wc.numberOfWorkers)
	for i := 0; i < wc.numberOfWorkers; i++ {
		go wc.urlWorker(i, urls)
	}
	wc.jobsQueue <- q
}

// Pick a url from the job queue, starts a worker to parse, puts this url to output
func (wc *WebCrawler) urlWorker(id int, urls chan<- UrlData) {
	log.Infof("At worker %d", id)
	for q := range wc.jobsQueue {
		if q == nil {
			wc.jobsQueue <- nil
			break
		}
		currentUrl := q.Dequeue()
		if currentUrl != nil {
			log.Infof("Processing at worker id:%d urlData:%s", id, currentUrl)
			go wc.parserWorker(currentUrl.(UrlData))

			go func() {
				urls <- currentUrl.(UrlData)
			}()
		}

		wc.jobsQueue <- q
	}
	log.Infof("Exiting worker %d", id)
}

// Picks a url to parse, puts new urls in the job queue, if not processed already
func (wc *WebCrawler) parserWorker(urlData UrlData) {
	newUrls := wc.siteParser.Parse(urlData.Url)
	log.Infof("Parsed url:%s result:%s", urlData, newUrls)
	q := <-wc.jobsQueue
	if q == nil {
		return
	}
	for _, newUrl := range newUrls {
		allUrls := q.Metadata[allUrls].(set.Interface)
		if !allUrls.Contains(newUrl) {
			allUrls.Add(newUrl)
			q.Enqueue(UrlData{
				Url:       newUrl,
				SourceUrl: urlData.Url,
			})
		}
	}

	wc.jobsQueue <- q
}

func (wc *WebCrawler) Stop() {
	wc.jobsQueue <- nil
}
