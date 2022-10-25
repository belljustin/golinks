package golinks

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/belljustin/golinks/pkg/golinks"
)

type HttpServer struct {
	service Service
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		service: defaultService(),
	}
}

func NewStorage() golinks.Storage {
	return golinks.NewStorage(C.Storage.Type)
}

func (s *HttpServer) Serve() error {
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

func (s *HttpServer) ping(w http.ResponseWriter, req *http.Request) {
	if err := s.service.Ping(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if _, err := fmt.Fprint(w, "pong"); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *HttpServer) home(w http.ResponseWriter, req *http.Request) {
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

func (s *HttpServer) links(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		s.postLink(w, req)
	} else {
		log.Printf("[INFO] method '%s' not allowed for path '%s'", req.Method, req.URL.Path)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *HttpServer) postLink(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Printf("[INFO] failed to create link: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := req.Form
	link, err := parseLinkValues(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	content, err := s.service.SetLink(link.Name, link.URL)
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