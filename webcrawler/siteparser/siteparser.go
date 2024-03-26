package siteparser

type Interface interface {
	Parse(string) ([]string, error)
}
