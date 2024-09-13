package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	chiFileServer "todo-list/internal/lib/chi-FileServer"
	"todo-list/internal/lib/logger"
	"todo-list/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const webPath = "../web"

func main() {

	//TODO init config

	//TODO init logger
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	log.Info("Starting TODO-list app.")

	//TODO init database
	storage, err := sqlite.NewStorage(log)
	if err != nil {
		log.Error("Failed to initialize database", logger.Err(err))
		return
	}
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/", http.FileServer(http.Dir(webPath)))

	log.Info("Configure fileserver.")
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, webPath))
	err = chiFileServer.FileServer(router, "/", filesDir)
	if err != nil {
		log.Error("Failed to сonfigure th fileserveer", logger.Err(err))
	}

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

// TODO Refactor this
