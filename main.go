package main

import (
	"go_final_project/function"
	"go_final_project/handler"
	"go_final_project/repository"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %v", err)
	}

	DBFile := os.Getenv("TODO_DBFILE")
	db := function.DBconnect()
	defer db.Close()
	repo := repository.NewDBwork(db)

	if _, err := os.Stat(DBFile); err != nil {
		if os.IsNotExist(err) {
			repo.AddDB(DBFile)
		} else {
			log.Fatalf("error while stat %s: %v", DBFile, err)
		}
	} else {
		err := repo.CheckTable()
		if err != nil {
			repo.CreateTable()
		}
	}
}
func main() {

	db := function.DBconnect()
	defer db.Close()

	repo := repository.NewDBwork(db)
	handlerfunc := handler.NewHandler(repo)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %v", err)
	}
	web_server_port := os.Getenv("TODO_PORT")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/task", handlerfunc.PostGetPutDeleteTask)
	mux.HandleFunc("/api/tasks", handlerfunc.GetTasksSearch)
	mux.HandleFunc("/api/task/done", handlerfunc.DoneTaskId)
	mux.HandleFunc("/api/nextdate", handlerfunc.NextData)
	err = http.ListenAndServe(web_server_port, mux)
	if err != nil {
		panic(err)
	}
}
