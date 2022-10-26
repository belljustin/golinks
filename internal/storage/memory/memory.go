package memory

import (
	"log"
	"net/url"
	"sync"

	"github.com/belljustin/golinks/pkg/golinks"
)

func init() {
	golinks.RegisterStorage("memory", newStorage)
}

type Storage struct {
	lock sync.RWMutex
	m    map[string]url.URL
}

func newStorage() golinks.Storage {
	return &Storage{
		lock: sync.RWMutex{},
		m:    make(map[string]url.URL),
	}
}

func (s *Storage) GetLink(name string) (*url.URL, error) {
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
	s.m[name] = l

	return nil
}

func (s *Storage) Health() error {
	return nil
}

func (s *Storage) Migrate() error {
	log.Println("[INFO] memory: nothing to migrate")
	return nil
}
