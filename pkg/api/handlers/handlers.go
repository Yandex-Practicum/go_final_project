package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"proj/pkg/api/service"
	"strconv"
	"time"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	time, err := time.Parse(service.LayoutDate, now)
	if err != nil {
		http.Error(w, "Ошибка парсинга полученного из пути параметра date", http.StatusBadRequest)
		return
	}

	res, err := service.NextDate(time, date, repeat)
	if err != nil {
		http.Error(w, "Ошибка вычисления следующей даты", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(res))

}

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

	task.ID = strconv.Itoa(int(id))

	type response struct {
		ID string `json:"id"`
	}
	var res response
	res.ID = task.ID

	jsByte, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Ошибка подготовки ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsByte)

}
