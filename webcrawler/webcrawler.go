package webcrawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/qulia/go-examples/webcrawler/siteparser"
	"github.com/qulia/go-qulia/lib/set"
	log "github.com/sirupsen/logrus"
)

type WebCrawler struct {
	// Synchronization
	lock      sync.Mutex
	done      chan bool
	closeOnce sync.Once
	inflight  sync.WaitGroup
	// Communication
	jobsQueue chan string
	// Props
	numberOfWorkers int
	maxUrlCount     int
	siteParser      siteparser.Interface
	timeout         time.Duration
	// State
	urls   set.Set[string]
	errors *compositeerror
	// Admission Control
	hostController *hostController
}

func NewWebCrawler(numWorkers, maxUrlCount int,
	timeout time.Duration,
	perHostRequestGapDuration time.Duration,
	siteParser siteparser.Interface) *WebCrawler {
	wc := WebCrawler{
		done:            make(chan bool),
		lock:            sync.Mutex{},
		inflight:        sync.WaitGroup{},
		jobsQueue:       make(chan string, 10),
		numberOfWorkers: numWorkers,
		maxUrlCount:     maxUrlCount,
		siteParser:      siteParser,
		timeout:         timeout,

		urls:           set.NewSet[string](),
		errors:         &compositeerror{},
		hostController: NewHostController(perHostRequestGapDuration),
	}
	return &wc
}

// Visits BFS manner starting at the startUrl and posts the urls
func (wc *WebCrawler) Visit(startUrl string) (set.Set[string], error) {
	// Start workers
	log.Infof("Starting %d workers...", wc.numberOfWorkers)
	for i := 0; i < wc.numberOfWorkers; i++ {
		wc.inflight.Add(1)
		go func(wid int) {
			wc.urlWorker(wid)
			wc.inflight.Done()
		}(i)
	}

	wc.jobsQueue <- startUrl
	t := time.AfterFunc(wc.timeout, func() {
		close(wc.done)
	})
	wc.inflight.Wait()
	t.Stop()
	return wc.urls, wc.errors
}

// Run multiple workers in parallel to process urls in the queue
// If throttled requeue
// If fails track error and drop url
func (wc *WebCrawler) urlWorker(id int) {
	log.Infof("At worker %d", id)
	defer log.Infof("Exiting worker %d", id)
	for {
		select {
		case <-wc.done:
			return
		case jurl := <-wc.jobsQueue:
			waitCh, ok, err := wc.hostController.Admit(jurl)
			if err != nil {
				wc.addError(fmt.Errorf("url dropped: %s", jurl))
				wc.addError(err)
				continue
			}

			if !ok {
				// Cannot continue, requeue
				wc.jobsQueue <- jurl
			} else {
				<-waitCh
				foundUrls, err := wc.siteParser.Parse(jurl)
				if err != nil {
					wc.addError(err)
					continue
				}
				for _, nu := range foundUrls {
					if wc.addUrl(nu) {
						wc.inflight.Add(1)
						go func(njurl string) {
							wc.jobsQueue <- njurl
							wc.inflight.Done()
						}(nu)
					}
				}
			}
		}
	}
}

func (wc *WebCrawler) addError(err error) {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	wc.errors.Add(err)
}

func (wc *WebCrawler) addUrl(u string) bool {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.urls.Len() == wc.maxUrlCount {
		wc.Stop()
		return false
	}

	if wc.urls.Contains(u) {
		return false
	}
	wc.urls.Add(u)
	return true
}

func (wc *WebCrawler) Stop() {
	wc.closeOnce.Do(func() {
		close(wc.done)
	})
}
