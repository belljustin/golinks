package golinks

import (
	"errors"
	"log"
	"net/url"
)

type Link struct {
	Name string
	URL  url.URL
}

func parseLinkValues(values url.Values) (*Link, error) {
	name := values.Get("linkName")
	if name == "" {
		log.Println("[INFO] missing param 'linkName'")
		return nil, errors.New("bad request: missing param 'linkName'")
	}

	sURL := values.Get("linkURL")
	if sURL == "" {
		log.Println("[INFO] missing param 'linkURL'")
		return nil, errors.New("bad request: missing param 'linkURL'")
	}

	URL, err := url.Parse(sURL)
	if err != nil {
		log.Printf("[INFO] param 'linkURL' does not contain a valid url: %s", err)
		return nil, errors.New("bad request: param 'linkURL' does not contain a valid url")
	}

	return &Link{
		Name: name,
		URL:  *URL,
	}, nil
}
