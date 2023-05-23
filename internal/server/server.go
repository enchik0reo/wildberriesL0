package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpSrv *http.Server
}

func New() *Server {
	s := &Server{}
	s.httpSrv = &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpSrv.Addr = ":" + port
	s.httpSrv.Handler = handler
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
