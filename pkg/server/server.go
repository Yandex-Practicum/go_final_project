package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"proj/pkg/api/handlers"
	"proj/pkg/api/middleware"
	"proj/pkg/api/service"

	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func ServerRun() {
	r := chi.NewRouter()

	r.Post("/api/signin", handlers.LoginHandle)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/api/nextdate", handlers.NextDayHandler)
		r.Post("/api/task", handlers.AddTaskHandler)       //
		r.Get("/api/tasks", handlers.GetTaskHandler)       //
		r.Get("/api/task", handlers.GetIDHandler)          //
		r.Put("/api/task", handlers.PutTaskHandler)        //
		r.Post("/api/task/done", handlers.DonePostHandler) //
		r.Delete("/api/task", handlers.DeleteTaskHandler)  //
	})

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	server := &http.Server{
		Addr:         ":" + service.PortGlobal(),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Запускаем сервер...")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Получен сигнал остановки сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер остановлен.")

}
