package chatbot

import (
	"fmt"
	"log"
	"net/http"
)

// server contains data for chatbot server
type server struct {
	cfg *Config
	srv *http.Server
}

// New returns a new chatbot server
func New(cfg Config) (*server) {
	c := &cfg
	
	muxHandler := http.NewServeMux()
	muxHandler.HandleFunc("/", c.chatHandler)
	
	srv := &http.Server{
		Handler: muxHandler,
		Addr:    fmt.Sprintf("%s:%d", c.ServerIP, c.Port),
	}
	
	return &server {
		cfg: c,
		srv: srv,
	}
}

// ListenAndServe wrapper for ListenAndServe function
func (s *server) ListenAndServe() error {
	log.Printf("http server listening at %v", s.srv.Addr)
	return s.srv.ListenAndServe()
}
