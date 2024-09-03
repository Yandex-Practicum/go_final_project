package steps

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("/api/nextdate", NextDate)
	mux.HandleFunc("/api/signin", auth)
	mux.HandleFunc("/api/task", authTask(selectFunc))
	mux.HandleFunc("/api/tasks/", authTask(searchHandler))
	mux.HandleFunc("/api/task/done", authTask(TaskDone))

	portStr, exists := os.LookupEnv("PORT")
	var currPort string
	if exists {
		currPort = portStr
	} else {
		currPort = "7540"
	}
	err := http.ListenAndServe((":" + currPort), mux)
	fmt.Printf("Прослушивание порта: %s", currPort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Завершение работы")

}

func selectFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskWM(w, r)
		return
	case http.MethodGet:
		GetTaskId(w, r)
		return
	case http.MethodPut:
		EditTask(w, r)
		return
	case http.MethodDelete:
		DeleteTask(w, r)
		return
	}
}
