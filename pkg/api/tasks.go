package api

import (
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Функция для проверки формата даты
func isValidDateFormat(date string) bool {
	_, err := time.Parse("02.01.2006", date)
	return err == nil
}

// Функция для преобразования даты из формата dd.mm.yyyy в yyyyMMdd
func formatDateForSearch(date string) string {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return ""
	}
	return parsedDate.Format("20060102")
}

// обработчик GET-запроса /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр search из строки запроса
	search := r.URL.Query().Get("search")

	// если в поисковом запросе есть дата, пытаемся её распарсить
	var dateSearch string
	if isValidDateFormat(search) {
		dateSearch = formatDateForSearch(search)
		search = "" // очищаем search, чтобы поиск по строкам не мешал
	}

	tasks, err := db.Tasks(50, search, dateSearch) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка получения задач: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	result := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, map[string]string{
			"id":      fmt.Sprint(t.ID),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})
	}
	writeJson(w, map[string]any{
		"tasks": result,
	}, http.StatusOK)
}
