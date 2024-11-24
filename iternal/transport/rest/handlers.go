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
		CallError("нет заголовка", w)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(TimeFormat)
	} else {
		date, err = time.Parse(TimeFormat, task.Date)
		if err != nil {
			CallError("неверный формат даты", w)
			return
		}
	}

	if now.After(date) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(TimeFormat)
		} else {
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				CallError("неверный формат", w)
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
		CallError("ошибка с базой данных", w)
		return
	}

	resp, err := json.Marshal(map[string]string{"id": strconv.Itoa(int(id))})
	if err != nil {
		CallError("не получилось создать напоминание", w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []services.Task{}

	count, err := database.GetCountOfTasks()
	if err != nil {
		CallError("ошибка с базой данных", w)
		return
	}

	if count > 0 {
		rows, err := database.GetAllTasks()
		if err != nil {
			CallError("ошибка с базой данных", w)
			return
		}
		for rows.Next() {
			var task services.Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				CallError("ошибка с базой данных", w)
				return
			}
			tasks = append(tasks, task)
		}
	}
	resp, err := json.Marshal(map[string]interface{}{
		"tasks": tasks,
	})
	if err != nil {
		CallError("ошибка сериализации данных", w)
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

func CallError(txt string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": txt})
}
