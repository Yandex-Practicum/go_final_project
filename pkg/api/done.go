package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func handleTaskDone(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем id из запроса
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "id задачи не указан"})
		return
	}

	// Получение задачи из БД
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Ошибка получения задачи из БД"})
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Задача не найдена"})
		return
	}

	// Удаляем задачу если нет повторов
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Не удалось удалить задачу"})
			return
		}
		json.NewEncoder(w).Encode(UpdateTaskResponce{})
		return
	}

	// Ищем следующую дату для повтора
	if task.Repeat != "" {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Не удалось получить следующую дату"})
			return
		}

		// Обновляем дату
		err = db.UpdateDate(next, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Не удалось обновить дату задачи"})
			return
		}
		json.NewEncoder(w).Encode(UpdateTaskResponce{})
	}
}
