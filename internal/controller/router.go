package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler *Handler, addr string) *Server {
	router := chi.NewRouter()

	router.Get("/devices/{guid}", handler.GetDeviceLogsByGuid)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) StartHttpServer() error {
	return s.httpServer.ListenAndServe()
}
