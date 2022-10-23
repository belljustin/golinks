package memory

import (
	"net/url"
	"sync"

	"github.com/belljustin/golinks"
)

type Storage struct {
	lock sync.RWMutex
	m    map[string]golinks.Link
}

func NewStorage() *Storage {
	return &Storage{
		lock: sync.RWMutex{},
		m:    make(map[string]golinks.Link),
	}
}

func (s *Storage) GetLink(name string) (*golinks.Link, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	link, ok := s.m[name]
	if !ok {
		return nil, nil
	}

	return &link, nil
}

func (s *Storage) SetLink(name string, l url.URL) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[name] = golinks.Link{
		Name: name,
		URL:  l,
	}

	return nil
}
