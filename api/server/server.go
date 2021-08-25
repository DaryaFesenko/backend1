package server

import (
	"context"
	"net/http"
	"time"
)

const (
	timeout        = 30
	timeoutContext = 2
)

type Server struct {
	srv http.Server
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       timeout * time.Second,
		WriteTimeout:      timeout * time.Second,
		ReadHeaderTimeout: timeout * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start() {
	go s.srv.ListenAndServe()
}
