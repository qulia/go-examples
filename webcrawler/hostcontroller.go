package webcrawler

import (
	"net/url"
	"sync"
	"time"

	"github.com/qulia/go-qulia/algo/ratelimiter"
	"github.com/qulia/go-qulia/algo/ratelimiter/leakybucket"
)

type hostController struct {
	m                         map[string]ratelimiter.RateLimiterBuffered
	perHostRequestGapDuration time.Duration
	lock                      sync.Mutex
}

func (hc *hostController) Admit(u string) (<-chan interface{}, bool, error) {
	hc.lock.Lock()
	defer hc.lock.Unlock()

	up, err := url.Parse(u)
	if err != nil {
		return nil, false, err
	}
	host := up.Host
	if hc.m[host] == nil {
		hc.m[host] = leakybucket.NewLeakyBucket(1, 1, hc.perHostRequestGapDuration)
	}

	waitCh, ok := hc.m[host].Allow()

	return waitCh, ok, nil
}

func NewHostController(perHostRequestGapDuration time.Duration) *hostController {
	return &hostController{m: make(map[string]ratelimiter.RateLimiterBuffered), perHostRequestGapDuration: perHostRequestGapDuration}
}
