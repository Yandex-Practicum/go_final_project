package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"main.go/base"
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

func nextDateHandle(w http.ResponseWriter, req *http.Request) {
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
		w.Write([]byte("wrong count value"))
		fmt.Println(err)
	}

	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nextDate, err := base.NextDate(now, date, repeat)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

func main() {
	fmt.Println("Запускаем сервер")
	//http.HandleFunc(`/`, mainHandle)
	// Устанавливаем обработчик для корневого URL
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc(`/api/nextdate`, nextDateHandle)
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
