package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ZhakiyaK/final_project/pkg/api"
)

type Server struct {
	logger *log.Logger
	port   string
}

const Port = 7540

func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger: logger,
		port:   strconv.Itoa(Port),
	}
}

func (s *Server) Start() error {
	api.Init()

	// Директория для статических файлов
	webDir := "./web"

	// Настройка обработчиков
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Логируем запуск сервера
	s.logger.Printf("Server starting on http://localhost:%s\n", s.port)
	s.logger.Printf("Serving static files from: %s", webDir)

	// Запуск сервера
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		s.logger.Printf("Server stopped with error: %v", err)
	} else {
		s.logger.Println("Server stopped gracefully")
	}
	return err
}
