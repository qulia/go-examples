package siteparser

type Interface interface {
	Parse(url string) []string
}
