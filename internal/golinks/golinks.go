package golinks

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/belljustin/golinks/pkg/golinks"
)

type Server struct {
	service Service
}

func NewServer() *Server {
	return &Server{
		service: defaultService(),
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
	if err := s.service.Ping(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if _, err := fmt.Fprint(w, "pong"); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		content, err := s.service.Home()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(content); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	name := strings.TrimLeft(req.URL.Path, "/")
	redirect, err := s.service.GetLink(name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if redirect == nil {
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, req, redirect.String(), http.StatusTemporaryRedirect)
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

	content, err := s.service.SetLink(name, *URL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(content); err != nil {
		log.Printf("[ERROR] failed to write: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
