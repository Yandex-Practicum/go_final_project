package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"todo-list/internal/http-server/handlers"
	chiFileServer "todo-list/internal/lib/chi-FileServer"
	"todo-list/internal/lib/logger"
	"todo-list/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const webPath = "../web"

func main() {

	//TODO init config
	// - databasePath
	// - http-server port
	// - project root folder

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	log.Info("Starting TODO-list app.")

	storage, err := sqlite.NewStorage(log)
	if err != nil {
		log.Error("Failed to initialize database", logger.Err(err))
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Handle("/", http.FileServer(http.Dir(webPath)))

	log.Debug("Configure fileserver.")
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, webPath))
	err = chiFileServer.FileServer(router, "/", filesDir)
	if err != nil {
		log.Error("Failed to сonfigure th fileserveer", logger.Err(err))
	}

	router.Get("/api/nextdate", handlers.GetNextDate(log))
	router.Post("/api/task", handlers.PostTask(log, storage))
	router.Get("/api/tasks", handlers.GetTasks(log, storage))
	router.Get("/api/task", handlers.GetTask(log, storage))
	router.Put("/api/task", handlers.PutTask(log, storage))
	router.Post("/api/task/done", handlers.MarkAsDone(log, storage))
	router.Delete("/api/task", handlers.DelTask(log, storage))

	server := http.Server{
		Addr:    "localhost:7540",
		Handler: router,
	}

	log.Info("Starting http-server")
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start http-server")
	}

	log.Error("Server stopped")
}
