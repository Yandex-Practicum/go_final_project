package api

import (
	"net/http"
	"time"

	"github.com/AndreiCoder404/go_final_project/pkg/db"
)

// TasksResp представляет структуру ответа для списка задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Извлекаем параметр search
	search := r.FormValue("search")
	var tasks []*db.Task
	var err error

	// Проверяем, является ли search датой в формате 02/01/2006
	if search != "" {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			// Если это дата, преобразуем в формат 20060102 и ищем задачи на эту дату
			dateStr := parsedDate.Format(DateFormat)
			tasks, err = db.Tasks(50, db.WithDate(dateStr))
		} else {
			// Иначе ищем как подстроку в title или comment
			tasks, err = db.Tasks(50, db.WithSearch(search))
		}
	} else {
		// Если search не указан, возвращаем все задачи
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve tasks")
		return
	}

	// Убеждаемся, что слайс не nil
	if tasks == nil {
		tasks = []*db.Task{}
	}

	// Формируем ответ
	resp := TasksResp{
		Tasks: tasks,
	}

	// Отправляем ответ
	writeJSON(w, resp)
}
