package golinks

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/belljustin/golinks/pkg/golinks"
)

type Server struct {
	storage golinks.Storage
	webPath string
}

func NewServer() *Server {
	storage := golinks.NewStorage(C.Storage.Type)

	return &Server{
		storage: storage,
		webPath: C.WebPath,
	}
}

func NewStorage() golinks.Storage {
	return golinks.NewStorage(C.Storage.Type)
}

func (s *Server) Serve() error {
	log.Printf("[INFO] starting golinks server on port :%s", C.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.ping)
	mux.HandleFunc("/links", s.links)
	mux.HandleFunc("/", s.home)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", C.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func (s *Server) ping(w http.ResponseWriter, req *http.Request) {
	if _, err := fmt.Fprintf(w, "pong"); err != nil {
		log.Printf("[ERROR] %s", err)
	}
}

func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		fName := path.Join(s.webPath, "html/index.html")
		http.ServeFile(w, req, fName)
		return
	}

	s.getLink(w, req)
}

func (s *Server) links(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		s.postLink(w, req)
	} else {
		log.Printf("[INFO] method '%s' not allowed for path '%s'", req.Method, req.URL.Path)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *Server) postLink(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Printf("[INFO] failed to create link: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := req.Form
	name := values.Get("linkName")
	if name == "" {
		log.Println("[INFO] missing param 'linkName'")
		http.Error(w, "Bad Request: missing param 'linkName'", http.StatusBadRequest)
		return
	}

	sURL := values.Get("linkURL")
	if sURL == "" {
		log.Println("[INFO] missing param 'linkURL'")
		http.Error(w, "Bad Request: missing param 'linkURL'", http.StatusBadRequest)
		return
	}
	URL, err := url.Parse(sURL)
	if err != nil {
		log.Printf("[INFO] param 'linkURL' does not contain a valid url: %s", err)
		http.Error(w, "Bad Request: param 'linkURL' does not contain a valid url", http.StatusBadRequest)
		return
	}

	if err := s.storage.SetLink(name, *URL); err != nil {
		log.Printf("[ERROR] failed to set link: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] link added")

	p := path.Join(s.webPath, "html/tmpl/link-added.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		log.Printf("[ERROR] failed to parse link-added.html: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  sURL,
	})
	if err != nil {
		log.Printf("[ERROR] failed to execute link-added.html template: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getLink(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimLeft(req.URL.Path, "/")

	link, err := s.storage.GetLink(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if link == nil {
		log.Printf("[INFO] link name '%s' does not exist", name)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("[INFO] redirect link name '%s' to '%s'", name, link.String())
	http.Redirect(w, req, link.String(), http.StatusTemporaryRedirect)
}
