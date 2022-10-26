package golinks

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/url"

	"github.com/belljustin/golinks/pkg/golinks"
)

const (
	indexFName         = "web/html/index.html"
	healthTmplFname    = "web/html/tmpl/health.html"
	linkAddedTmplFName = "web/html/tmpl/link-added.html"
)

//go:embed web
var resources embed.FS

func init() {
	// verify the static files have been properly embedded
	if _, err := resources.ReadFile(indexFName); err != nil {
		panic(err)
	}
	if _, err := resources.ReadFile(linkAddedTmplFName); err != nil {
		panic(err)
	}
}

type HealthCheck struct {
	Component string
	Error     error
}

type HealthChecks []HealthCheck

func (hs HealthChecks) HTML() ([]byte, error) {
	t, err := template.ParseFS(resources, healthTmplFname)
	if err != nil {
		log.Printf("[ERROR] failed to parse %s: %s", healthTmplFname, err)
		return nil, err
	}

	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, hs)
	if err != nil {
		log.Printf("[ERROR] failed to execute %s template: %s", healthTmplFname, err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

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
	Home() ([]byte, error)
	GetLink(name string) (*url.URL, error)
	SetLink(name string, link url.URL) ([]byte, error)
}

type serviceImpl struct {
	storage golinks.Storage
}

func defaultService() Service {
	storage := golinks.NewStorage(C.Storage.Type)

	return &serviceImpl{
		storage: storage,
	}
}

func (svc *serviceImpl) Health() HealthChecks {
	return HealthChecks{
		{"server", nil},
		{"storage", svc.storage.Health()},
	}
}

func (svc *serviceImpl) Home() ([]byte, error) {
	content, err := resources.ReadFile(indexFName)
	if err != nil {
		log.Printf("[ERROR] failed to read '%s': %s", indexFName, err)
	}
	return content, err
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

func (svc *serviceImpl) SetLink(name string, link url.URL) ([]byte, error) {
	if err := svc.storage.SetLink(name, link); err != nil {
		log.Printf("[ERROR] failed to set link: %s", err)
		return nil, err
	}
	log.Println("[INFO] link added")

	t, err := template.ParseFS(resources, linkAddedTmplFName)
	if err != nil {
		log.Printf("[ERROR] failed to parse %s: %s", linkAddedTmplFName, err)
		return nil, err
	}

	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  link.String(),
	})
	if err != nil {
		log.Printf("[ERROR] failed to execute link-added.html template: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
