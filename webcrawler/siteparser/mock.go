package siteparser

import "time"

type MockSiteParser struct {
	reachableUrlMap map[string][]string
}

func NewMockSiteParser(urlMap map[string][]string) Interface {
	msp := MockSiteParser{reachableUrlMap: urlMap}
	return &msp
}

func (msp *MockSiteParser) Parse(url string) []string {
	if val, ok := msp.reachableUrlMap[url]; !ok {
		return nil
	} else {
		time.Sleep(time.Microsecond)
		return val
	}
}
