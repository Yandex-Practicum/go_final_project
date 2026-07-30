package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"proj/pkg/api/service"

	"strconv"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task service.Task

	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		service.SendErrorJSON(w, "Ошибка чтения body из запроса", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(byteBody, &task)
	if err != nil {
		service.SendErrorJSON(w, "Ошибка дисериализации json из прочитанного body", http.StatusInternalServerError)
		return
	}
	if task.Title == "" {
		service.SendErrorJSON(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	err = service.CheckDate(&task)
	if err != nil {
		service.SendErrorJSON(w, "Ошибка при проверке параметров.", http.StatusBadRequest)
		return
	}

	id, err := service.AddTask(&task)
	if err != nil {
		service.SendErrorJSON(w, "Ошибка добавления новой задачи", http.StatusInternalServerError)
		return
	}

	// task.ID = strconv.Itoa(int(id))

	type response struct {
		ID string `json:"id"`
	}
	var res response

	res.ID = strconv.Itoa(int(id))

	jsByte, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Ошибка подготовки ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsByte)

}
