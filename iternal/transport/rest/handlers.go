package handlers

import (
	"Go/iternal/database"
	"Go/iternal/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	TimeFormat = "20060102"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task services.Task
	var buf bytes.Buffer
	var date time.Time

	now, _ := time.Parse(TimeFormat, time.Now().Format(TimeFormat))

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "нет заголовка"})
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(TimeFormat)
	} else {
		date, err = time.Parse(TimeFormat, task.Date)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(map[string]string{"error": "неверный формат даты"})
			return
		}
	}

	if now.After(date) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(TimeFormat)
		} else {
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				json.NewEncoder(w).Encode(map[string]string{"error": "неверный формат"})
				return
			}
		}
	}
	//else if task.Repeat != "" {
	//	task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
	//	if err != nil {
	//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//		json.NewEncoder(w).Encode(map[string]string{"error": "неверный формат"})
	//		return
	//	}
	// }

	id, err := database.PutTaskInDB(task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "ошибка с базой данных"})
		return
	}

	resp, err := json.Marshal(map[string]string{"id": strconv.Itoa(int(id))})
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "не получилось создать напоминание"})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func NextDeadLine(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(TimeFormat, r.URL.Query().Get("now"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	deadline, err := services.NextDate(now, date, repeat)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Write([]byte(deadline))

}
