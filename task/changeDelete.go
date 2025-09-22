package task

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// UpdateTasks обрабатывает PUT-запросы для редактирования задачи  6ой шаг
func UpdateTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)
	if req.Method != http.MethodPut {
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}

	var task Task
	// Декодируем JSON из запроса
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&task); err != nil {
		log.Printf("Ошибка десериализации JSON")
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка на наличие идентификатора
	if task.ID == "" {
		log.Printf("Не указан идентификатор")
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	// Проверка, что идентификатор является числом
	_, err := strconv.Atoi(task.ID)
	if err != nil {
		log.Printf("Некорректный идентификатор")
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}
	// Проверяем обязательное поле title
	if task.Title == "" {
		log.Printf("Не указан заголовок задачи")
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверяем парсится ли дата
	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		log.Printf("Неправильный формат даты")
		http.Error(w, `{"error": "Неправильный формат даты"}`, http.StatusBadRequest)
		return
	}

	if taskDate.Before(time.Now().Local().Truncate(24 * time.Hour)) {
		log.Printf("Дата не может быть в прошлом")
		http.Error(w, `{"error": "Дата не может быть в прошлом"}`, http.StatusBadRequest)
		return
	}
	// Проверяем обязательное поле repeat
	if task.Repeat == "" {
		log.Printf("Неправильно указано правило повторения")
		http.Error(w, `{"error": "Неправильно указано правило повторения"}`, http.StatusBadRequest)
		return
	}
	if !regexp.MustCompile(`^(d \d{1,3}|y)$`).MatchString(task.Repeat) {
		log.Printf("Неправильно указано правило повторения")
		http.Error(w, `{"error": "Неправильно указано правило повторения"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли задача с указанным id
	var exist bool
	err = db.QueryRow("SELECT 1 FROM scheduler WHERE id = ?", task.ID).Scan(&exist)
	if err != nil {
		log.Printf("Не удалось сделать запрос к БД")
		http.Error(w, `{"error": "Не удалось сделать запрос к БД"}`, http.StatusNotFound)
		return
	}
	if !exist {
		log.Printf("Задача не найдена")
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Обновляем задачу в БД
	_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("Не удалось обновить задачу")
		http.Error(w, `{"error": "Не удалось обновить задачу"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// DoneTask обрабатывает PUT-запросы для выполненной задачи 7ой шаг
func DoneTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

	id := req.URL.Query().Get("id")
	if id == "" {
		log.Printf("Не указан идентификатор задачи")
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error": "Неверный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	var t Task
	row := db.QueryRow("SELECT  date, repeat FROM scheduler WHERE id = ?", id)
	err = row.Scan(&t.Date, &t.Repeat)
	if err == sql.ErrNoRows {
		log.Printf("Задача не найдена в БД")
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Ошибка запроса к БД")
		http.Error(w, `{"error": "Ошибка запроса к БД"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Запрашиваем задачу в DoneTask с ID: %s", id)

	if t.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Не удалось удалить задачу")
			http.Error(w, `{"error": "Не удалось удалить задачу"}`, http.StatusNotFound)
			return
		}
		// Возвращаем пустой JSON в случае успеха
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{}"))
		return
	}

	newDate, err := NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		log.Printf("Не удалось вычислить следующую дату")
		http.Error(w, `{"error": "Не удалось вычислить следующую дату"}`, http.StatusNotFound)
		return
	}

	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
	if err != nil {
		http.Error(w, `{"error": "Не удалось обновить задачу"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

// DeleteTask обрабатывает DELETE-запросы для удаления задачи 7ой шаг
func DeleteTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)
	if req.Method != http.MethodDelete {
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}
	id := req.URL.Query().Get("id")
	if id == "" {
		log.Printf("Не указан идентификатор задачи")
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Удаляем задачу из БД
	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		log.Printf("Не удалось удалить задачу")
		http.Error(w, `{"error": "Не удалось удалить задачу"}`, http.StatusNotFound)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
