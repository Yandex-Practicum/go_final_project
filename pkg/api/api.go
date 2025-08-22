package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const DateFormat = "20060102"

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	gnow := req.URL.Query().Get("now")
	if gnow == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("now missing"))
		return
	}
	parsTime, err := time.Parse(DateFormat, gnow)
	if err != nil {
		fmt.Println("Ошибка при разборе даты dstart", err)
		return
	}

	gdate := req.URL.Query().Get("date")
	if gdate == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("date missing"))
		return
	}

	grepeat := req.URL.Query().Get("repeat")
	if grepeat == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("repeat missing"))
		return
	}

	str, err := NextDate(parsTime, gdate, grepeat)
	if err != nil {
		fmt.Println("Ошибка NextDate:", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("NextDate missing"))
		return
	}

	res.Write([]byte(str))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}

func selectTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	}
}

func onetaskkHandler(w http.ResponseWriter, r *http.Request) {
	getoneTaskHandler(w, r)
}

func putttaskkHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		putTaskHandler(w, r)
	}
}

func deletetaskkHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		delTaskHandler(w, r)
	}
}

func donehandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		doneTaskHandler(w, r)
	}
}

func Init() *chi.Mux {
	rout := chi.NewRouter()
	rout.Handle("/*", http.FileServer(http.Dir("web")))
	rout.Get("/api/nextdate", nextDayHandler)
	rout.Post("/api/task", taskHandler)
	rout.Get("/api/tasks", selectTaskHandler)
	rout.Get("/api/task", onetaskkHandler)
	rout.Put("/api/task", putttaskkHandler)
	rout.Delete("/api/task", deletetaskkHandler)
	rout.Post("/api/task/done", donehandler)

	return rout
}
