package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"todo/config"
	httpserver "todo/infrastructure/http_server"
	"todo/infrastructure/jwt"
	sqliterepository "todo/infrastructure/sqlite_repository"
	"todo/internal/api/http"
	"todo/internal/api/http/controller"
	"todo/internal/usecase"
)

func Run(cfg *config.Config) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	repository, err := sqliterepository.New(cfg.DBFile)
	if err != nil {
		return fmt.Errorf("sqliterepository.New: %w", err)
	}

	defer repository.Close()

	controller := controller.NewTaskController(
		usecase.NewTaskUsecase(repository, jwt.New(cfg.JWTSecret), usecase.WithPassword(cfg.Password)))
	router := http.NewRouter(controller)
	server := httpserver.New(router.Handler(), httpserver.WithPort(cfg.Port))
	server.Start()
	log.Println("Server started")

	select {
	case <-quit:
		if err = server.Shutdown(); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	case err = <-server.Notify():
		return fmt.Errorf("server.Notify: %w", err)
	}
}
