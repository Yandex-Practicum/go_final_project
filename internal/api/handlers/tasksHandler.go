package handlers

import (
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"
	"go-task-manager-final_project/internal/services"

	"net/http"
	"strings"
	"time"
)

// TasksResp - структура для ответа API, содержит список задач.
// Поле Tasks представляет собой слайс указателей на задачи из БД.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler - обработчик HTTP-запросов для получения списка задач.
// Поддерживает фильтрацию по поисковому запросу (поиск по заголовку, комментарию или дате).
// Параметры:
// w - объект для записи HTTP-ответа;
// r - объект HTTP-запроса.
func (s *APIServer) tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр search из строки запроса
	searchQuery := r.URL.Query().Get("search")

	// Вызываем БД для получения списка задач (максимум 50 записей)
	tasks, err := db.GetTasks(s.DB, 50)
	if err != nil {
		// Возвращаем HTTP 500 с сообщением об ошибке
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch tasks from database",
		})
		return
	}

	// Если задач нет - возвращаем пустой массив, а не null
	if tasks == nil {
		tasks = []*db.Task{}
	}

	// Если есть поисковый запрос - фильтруем задачи
	if searchQuery != "" {
		filteredTasks := []*db.Task{}

		// Проверяем, является ли searchQuery датой в формате services.DateFormat
		isDate := false
		parsedDate, err := time.Parse(services.DateFormat, searchQuery)
		if err == nil {
			isDate = true
		}

		// Если не получилось, пробуем альтернативный формат DD.MM.YYYY
		if !isDate {
			parsedDate, err = time.Parse("02.01.2006", searchQuery)
			isDate = err == nil
		}

		// Проходим по всем задачам и отбираем подходящие под фильтр
		for _, task := range tasks {
			if isDate {
				// Преобразуем строку из задачи в time.Time
				taskDate, err := time.Parse(services.DateFormat, task.Date)
				if err != nil {
					taskDate, err = time.Parse("02.01.2006", task.Date)
					if err != nil {
						continue
					}
				}
				// Сравниваем даты на равенство
				if taskDate.Equal(parsedDate) {
					filteredTasks = append(filteredTasks, task)
				}
			} else {
				// Проверяем, содержится ли поисковая строка в заголовке или комментарии (без учёта регистра)
				if strings.Contains(strings.ToLower(task.Title), strings.ToLower(searchQuery)) || strings.Contains(strings.ToLower(task.Comment), strings.ToLower(searchQuery)) {
					filteredTasks = append(filteredTasks, task)
				}
			}
		}
		tasks = filteredTasks
	}

	// Формируем и отправляем ответ в формате JSON с кодом 200 (OK)
	api.WriteJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
