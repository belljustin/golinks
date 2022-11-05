package golinks

import (
	"log"
	"net/url"

	"github.com/belljustin/golinks/pkg/golinks"
)

type serviceImpl struct {
	storage golinks.Storage
}

func defaultService() golinks.Service {
	storage := golinks.NewStorage(C.Storage.Type)

	return &serviceImpl{
		storage: storage,
	}
}

func (svc *serviceImpl) Health() golinks.HealthChecks {
	return golinks.HealthChecks{
		{"server", nil},
		{"storage", svc.storage.Health()},
	}
}

func (svc *serviceImpl) GetLink(name string) (*url.URL, error) {
	link, err := svc.storage.GetLink(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil, err
	}

	if link == nil {
		log.Printf("[INFO] link name '%s' does not exist", name)
		return nil, nil
	}

	log.Printf("[INFO] redirect link name '%s' to '%s'", name, link.String())
	return link, nil
}

func (svc *serviceImpl) SetLink(name string, link url.URL) error {
	if err := svc.storage.SetLink(name, link); err != nil {
		log.Printf("[ERROR] failed to set link: %s", err)
		return err
	}
	log.Println("[INFO] link added")
	return nil
}
