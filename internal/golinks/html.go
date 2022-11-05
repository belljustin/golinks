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

type htmlMarshaller struct {
}

func (m *htmlMarshaller) healthChecks(hs golinks.HealthChecks) ([]byte, error) {
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

func (m *htmlMarshaller) home() ([]byte, error) {
	content, err := resources.ReadFile(indexFName)
	if err != nil {
		log.Printf("[ERROR] failed to read '%s': %s", indexFName, err)
	}
	return content, err
}

func (m *htmlMarshaller) setLink(name string, link url.URL) ([]byte, error) {
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
