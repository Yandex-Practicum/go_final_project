package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"final_project/internal/api"
	"final_project/internal/api/tasks"
	"final_project/tests"
)

func CreateServer() {
	logger := log.New(os.Stdout, "[server]", log.LstdFlags|log.Lshortfile)
	todoPort := os.Getenv("TODO_PORT")
	if todoPort == "" {
		todoPort = strconv.Itoa(tests.Port)
	}
	router := http.NewServeMux()
	webDir := "./web"
	router.Handle("/", http.FileServer(http.Dir(webDir)))
	router.HandleFunc("GET /api/nextdate", api.NextDayHandler)
	router.HandleFunc("POST /api/task", tasks.AddTaskHandler)
	router.HandleFunc("GET /api/tasks", tasks.GetTasksHandler)
	router.HandleFunc("GET /api/task", tasks.GetTaskHandler)
	router.HandleFunc("PUT /api/task", tasks.UpdateTaskHandler)
	router.HandleFunc("DELETE /api/task", tasks.DeleteTaskHandler)
	router.HandleFunc("POST /api/task/done", tasks.TaskDoneHandler)
	server := &http.Server{
		Addr:         ":" + todoPort,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal("Ошибка при запуске сервера: ", err)
	}
}
