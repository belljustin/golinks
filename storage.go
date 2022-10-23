package golinks

import "net/url"

type Link struct {
	Name string
	URL  url.URL
}

type Storage interface {
	GetLink(name string) (*Link, error)
	SetLink(name string, url url.URL) error
}
