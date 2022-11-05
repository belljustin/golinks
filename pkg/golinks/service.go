package golinks

import (
	"net/url"
)

type HealthCheck struct {
	Component string
	Error     error
}

type HealthChecks []HealthCheck

func (hs HealthChecks) Error() bool {
	for _, h := range hs {
		if h.Error != nil {
			return true
		}
	}
	return false
}

type Service interface {
	Health() HealthChecks
	GetLink(name string) (*url.URL, error)
	SetLink(name string, link url.URL) error
}
