package addtask

import (
	"encoding/json"
	"final_project/internal/common"
	nextdate "final_project/internal/handlers/next_date"
	"net/http"
	"time"
)

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "application/json")
	tsk := common.Task{}
	err := json.NewDecoder(r.Body).Decode(&tsk)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "ошибка обработки запроса"})
		return
	}
	r.Body.Close()

	if tsk.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указан заголовок залачи"})
		return
	}
	if tsk.Date == "" {
		tsk.Date = time.Now().Format(common.TimeFormat)

	}

	dateTime, err := time.Parse(common.TimeFormat, tsk.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "неверный формат даты отсчёта"})
		return
	}
	dateTime = dateTime.Truncate(24 * time.Hour)
	now := time.Now().Truncate(24 * time.Hour)

	if dateTime.Before(now) {
		if tsk.Repeat == "" {
			dateTime = time.Now()
			tsk.Date = dateTime.Format(common.TimeFormat)
		} else {
			tsk.Date, err = nextdate.NextDate(time.Now(), tsk.Date, tsk.Repeat)
			if err != nil {
				json.NewEncoder(w).Encode(common.Response{Error: err.Error()})
				return
			}
		}
	}
	id, er := h.repo.AddTask(tsk.Date, tsk.Title, tsk.Comment, tsk.Repeat)
	if er != "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: er})
		return
	}
	json.NewEncoder(w).Encode(common.Response{ID: id})
}
