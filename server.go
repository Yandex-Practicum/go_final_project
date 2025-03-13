package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ServerStart() {

	if len(os.Getenv("TODO_PORT")) > 0 {
		port = os.Getenv("TODO_PORT")
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Getenv("TODO_DBFILE")) > 0 {
		DBFile = filepath.Join(filepath.Dir(appPath), "TODO_DBFILE")
	}

	_, err = os.Stat(DBFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		CreateDB(DBFile)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)

	http.HandleFunc("/api/tasks", getTasksHandler)

	http.HandleFunc("/api/nextdate", nextDateHandler)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
	return
}
