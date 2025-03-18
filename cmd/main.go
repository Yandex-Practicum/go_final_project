package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/wissio/go_final_project/internal/config"
	"github.com/wissio/go_final_project/internal/http-server/handlers"
	"github.com/wissio/go_final_project/internal/lib/logger"
	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	registerRoutes(router, log, storage)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Info("Starting server", slog.String("address", cfg.Address))
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
	}
}

func registerRoutes(router *chi.Mux, loggerInstance *slog.Logger, storage *sqlite.Storage) {
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	registerAPIRoutes(router, loggerInstance, storage)
}

func registerAPIRoutes(router *chi.Mux, loggerInstance *slog.Logger, storage *sqlite.Storage) {
	router.Get("/api/nextdate", handlers.NextDate(loggerInstance, storage))
	router.Post("/api/task", handlers.CreateTask(loggerInstance, storage))
	router.Get("/api/task", handlers.GetTask(loggerInstance, storage))
	router.Get("/api/tasks", handlers.GetTasks(loggerInstance, storage))
	router.Put("/api/task", handlers.UpdateTask(loggerInstance, storage))
	router.Post("/api/task/done", handlers.DoneTask(loggerInstance, storage))
	router.Delete("/api/task", handlers.DeleteTask(loggerInstance, storage))
}
