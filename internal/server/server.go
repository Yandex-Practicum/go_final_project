package server

import (
	"go_final_project/internal/config"
	_ "go_final_project/internal/config"
	"log"
	"net/http"
)

type Server struct {
  httpServer *http.Server
  Handler http.Handler
}

var port = config.Port


func (s *Server) Run(r http.Handler) error {
  s.httpServer = &http.Server{
    Addr: ":"+port,
    Handler: r,
  }

  log.Printf("Запуск сервера на %s", s.httpServer.Addr)

  return s.httpServer.ListenAndServe()
}
