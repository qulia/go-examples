package webcrawler

import (
	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/concurrency/access"
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
	jobsQueue       *access.Unique
}

func NewWebCrawler(numWorkers int, siteParser siteparser.Interface) *WebCrawler {
	wc := WebCrawler{
		numberOfWorkers: numWorkers,
		siteParser:      siteParser,
		jobsQueue:       access.NewUnique(queue.NewQueue()),
	}

	wc.jobsQueue.Release()
	return &wc
}

// Visits BFS manner starting at the startUrl and posts the urls
func (wc *WebCrawler) Visit(startUrl string, urls chan<- UrlData) {
	q := wc.jobsQueue.Acquire().(*queue.Queue)
	defer wc.jobsQueue.Release()

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
}

// Pick a url from the job queue, starts a worker to parse, puts this url to output
func (wc *WebCrawler) urlWorker(id int, urls chan<- UrlData) {
	log.Infof("At worker %d", id)
	for {
		q := wc.jobsQueue.Acquire()
		if q == nil {
			wc.jobsQueue.Release()
			break
		}
		currentUrl := q.(*queue.Queue).Dequeue()
		if currentUrl != nil {
			log.Infof("Processing at worker id:%d urlData:%s", id, currentUrl)
			go wc.parserWorker(currentUrl.(UrlData))

			go func() {
				urls <- currentUrl.(UrlData)
			}()
		}

		wc.jobsQueue.Release()
	}
	log.Infof("Exiting worker %d", id)
}

// Picks a url to parse, puts new urls in the job queue, if not processed already
func (wc *WebCrawler) parserWorker(urlData UrlData) {
	newUrls := wc.siteParser.Parse(urlData.Url)
	log.Infof("Parsed url:%s result:%s", urlData, newUrls)
	q := wc.jobsQueue.Acquire()
	defer wc.jobsQueue.Release()
	if q == nil {
		return
	}
	for _, newUrl := range newUrls {
		allUrls := q.(*queue.Queue).Metadata[allUrls].(set.Interface)
		if !allUrls.Contains(newUrl) {
			allUrls.Add(newUrl)
			q.(*queue.Queue).Enqueue(UrlData{
				Url:       newUrl,
				SourceUrl: urlData.Url,
			})
		}
	}
}

func (wc *WebCrawler) Stop() {
	wc.jobsQueue.Done()
}
