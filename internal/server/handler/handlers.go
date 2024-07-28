package handler

import (
	"encoding/json"
	"go_final_project/internal/config"
	"go_final_project/internal/storage"
	"go_final_project/internal/task"
	"log"
	"net/http"
	"strconv"
	"time"
)

var webDir = config.WebDir

func GetFront() http.Handler {
	return http.FileServer(http.Dir(webDir))
}

func NextDate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, "Некоректная дата", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	res, err := task.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "string")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(res))

	if err != nil {
		log.Println("Ошибка записи: ", err)
	}
}

func AddTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Ошибка дессирилиализации JSON", http.StatusBadRequest)
			return
		}

		err = task.Check(&t)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		id, err := db.AddTask(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := map[string]string{
			"id": strconv.Itoa(id),
		}

		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(resp)
	}
}

func GetTasks(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []task.Task
		var err error

		tasks, err = db.GetList()

		if err != nil {
			log.Fatalf("Не удалось получить список задач %v", err)
		}

		if len(tasks) == 0 {
			tasks = []task.Task{}
		}

		res := map[string][]task.Task{
			"tasks": tasks,
		}

		resp, err := json.Marshal(res)
		if err != nil {
			log.Printf("ошибка десериализации: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("не удалось записать: %v", err)
		} else {
			log.Printf("Успешный вывод задач. Задач найденно %d", len(tasks))
		}
	}
}

func GetTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			log.Println("Не указан идентификатор")
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
			return
		}

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Не корректный id, не является числом")
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Не корректный id"})
			return
		}

		task, err := db.GetTask(id)
		if err != nil {
			log.Println("Не удалось получить задачу", err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := json.Marshal(task)
		if err != nil {
			log.Println("Не удалось десерилиазивать задачу: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(res)
		if err != nil {
			log.Println("Не удалось изменить задачу", err)
		}
	}
}

func ChangeTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error": Не удалось прочитать ответ}`, http.StatusBadRequest)
			return
		}

		err = task.Check(&t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		err = db.ChangeTask(t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		t, err = db.GetTask(t.Id)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func DoneTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id не является числом", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "id не является числом"})
			return
		}

		t, err := db.GetTask(id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не удалось получить задачу по id: " + id})
			return
		}

		if t.Repeat == "" {
			err = db.DeleteTask(id)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "не удалось удалить задачу по id: " + id})
				return
			}
		}

		if t.Repeat != "" {
			t.Date, err = task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				return
			}

			err = db.ChangeTask(t)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				return
			}
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(map[string]string{})
		if err != nil {
			log.Println("Не удалось записать ответ", err)
		}
	}
}

func DeleteTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id не является числом", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "id не является числом"})
			return
		}

		err = db.DeleteTask(id)
		if err != nil {
			log.Println("Не удачное удаление задачи")
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(map[string]string{})
		if err != nil {
			log.Println("err encode:", err)
			http.Error(w, `{"error":"Не удобное декодирование ответа"}`, http.StatusInternalServerError)
		}
	}
}
