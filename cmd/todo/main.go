package main

import (
	addtask "cactus3d/go_final_project/internal/http-server/handlers/add-task"
	deletetask "cactus3d/go_final_project/internal/http-server/handlers/delete-task"
	donetask "cactus3d/go_final_project/internal/http-server/handlers/done-task"
	gettask "cactus3d/go_final_project/internal/http-server/handlers/get-task"
	gettasks "cactus3d/go_final_project/internal/http-server/handlers/get-tasks"
	nextdate "cactus3d/go_final_project/internal/http-server/handlers/next-date"
	updatetask "cactus3d/go_final_project/internal/http-server/handlers/update-task"
	"cactus3d/go_final_project/internal/service/tasks"
	"cactus3d/go_final_project/internal/storage/sqlite"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	WebDir = "./web/"
	DBFile = "./scheduler.db"
	Port   = 7540
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = DBFile
	}
	store, err := sqlite.New(dbFile)
	if err != nil {
		log.Fatalf("error starting db: %v", err)
		return
	}

	portStr := os.Getenv("TODO_PORT")
	var port int

	if portStr == "" {
		port = Port
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
			return
		}
	}

	taskService := tasks.New(store)

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir(WebDir)))

	r.MethodFunc(http.MethodGet, "/api/nextdate", nextdate.New())

	r.MethodFunc(http.MethodPost, "/api/task", addtask.New(taskService))
	r.MethodFunc(http.MethodGet, "/api/task", gettask.New(taskService))
	r.MethodFunc(http.MethodPut, "/api/task", updatetask.New(taskService))
	r.MethodFunc(http.MethodDelete, "/api/task", deletetask.New(taskService))

	r.MethodFunc(http.MethodGet, "/api/tasks", gettasks.New(taskService))

	r.MethodFunc(http.MethodPost, "/api/task/done", donetask.New(taskService))

	address := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
		return
	}

}
