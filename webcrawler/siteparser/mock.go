package siteparser

import (
	"fmt"
	"time"
)

type MockSiteParser struct {
	reachableUrlMap map[string][]string
	errorUrlMap     map[string]bool
}

func NewMockSiteParser(urlMap map[string][]string, errorMap map[string]bool) Interface {
	msp := MockSiteParser{reachableUrlMap: urlMap, errorUrlMap: errorMap}
	return &msp
}

func (msp *MockSiteParser) Parse(url string) ([]string, error) {
	if msp.errorUrlMap[url] {
		return nil, fmt.Errorf("could not pars url %s", url)
	}

	if val, ok := msp.reachableUrlMap[url]; !ok {
		return nil, nil
	} else {
		time.Sleep(time.Microsecond)
		return val, nil
	}
}
