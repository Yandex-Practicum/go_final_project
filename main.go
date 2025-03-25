package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"main.go/base"
	"main.go/task"
)

/* //еще один возможный вариант
func mainHandle(res http.ResponseWriter, req *http.Request) {
	var filePath string

	if req.URL.Path == "/" {
		filePath = filepath.Join("web", "index.html")
	} else {
		filePath = filepath.Join("web", req.URL.Path)
	}
	http.ServeFile(res, req, filePath)

}*/

const webDir = "./web"

func nextDateHandler(w http.ResponseWriter, req *http.Request) {
	nowStr := req.URL.Query().Get("now")
	if nowStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("now missing"))
		err := errors.New("пропущено время")
		fmt.Println(err)
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong time value"))
		fmt.Println(err)
	}

	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nextDate, err := task.NextDate(now, date, repeat)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

func taskHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	//case http.MethodGet: getTask(w, req)
	case http.MethodPost:
		task.AddTask(w, req)
	default:
		http.Error(w, fmt.Sprintf("Сервер не поддерживает %s запросы", req.Method),
			http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	fmt.Println("Запускаем сервер")
	//http.HandleFunc(`/`, mainHandle)
	// Устанавливаем обработчик для корневого URL
	/*db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return

	}
	defer db.Close()*/

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc(`/api/nextdate`, nextDateHandler)
	http.HandleFunc(`/api/task`, taskHandler)
	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		envPort = "7540"
	}
	envDBFILE := os.Getenv("TODO_DBFILE")
	base.CreateDB(envDBFILE)
	err := http.ListenAndServe(":"+envPort, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}
