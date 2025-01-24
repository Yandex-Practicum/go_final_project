package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	utils "github.com/falsefood/go_final_project/internal"
)

func UpdateTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("--- Функция updateTaskHandler вызвана ---")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела запроса: %v\n", err)
		utils.WriteError(w, http.StatusBadRequest, "Ошибка чтения тела запроса")
		return
	}
	defer r.Body.Close()

	log.Printf("Тело запроса (после чтения): %s\n", string(body))

	if len(body) == 0 {
		log.Println("Тело запроса пустое")
		utils.WriteError(w, http.StatusBadRequest, "Тело запроса не может быть пустым")
		return
	}

	input := utils.Task{}

	if err := json.Unmarshal(body, &input); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		utils.WriteError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	log.Printf("Декодированные данные: ID=%s, Date=%s, Title=%s, Comment=%s, Repeat=%s\n", input.ID, input.Date, input.Title, input.Comment, input.Repeat)

	id, err := strconv.Atoi(input.ID)
	if err != nil || id == 0 {
		log.Println("Ошибка: ID задачи не указан")
		utils.WriteError(w, http.StatusBadRequest, "ID задачи не указан")
		return
	}

	if input.Title == "" {
		log.Println("Ошибка: Заголовок задачи не указан")
		utils.WriteError(w, http.StatusBadRequest, "Заголовок нужно обязательно указать")
		return
	}

	today := time.Now()
	if input.Date == "" || input.Date == "today" {
		input.Date = today.Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", input.Date)
		if err != nil {
			log.Printf("Ошибка: Неверный формат даты: %v\n", err)
			utils.WriteError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}

		if parsedDate.Before(today) {
			if input.Repeat == "" {
				input.Date = today.Format("20060102")
			} else {
				nextDate, err := nextDate(today, input.Date, input.Repeat)
				if err != nil {
					log.Printf("Ошибка: Не удалось вычислить следующую дату: %v\n", err)
					utils.WriteError(w, http.StatusBadRequest, err.Error())
					return
				}
				input.Date = nextDate
			}
		}
	}

	if input.Repeat != "" {
		if !strings.HasPrefix(input.Repeat, "d ") && !strings.HasPrefix(input.Repeat, "y ") {
			log.Printf("Ошибка: Неверный формат повторения: %s\n", input.Repeat)
			utils.WriteError(w, http.StatusBadRequest, "Неверный формат повторения")
			return
		}
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	log.Printf("Выполнение SQL-запроса: %s с параметрами: Date=%s, Title=%s, Comment=%s, Repeat=%s, ID=%s\n", query, input.Date, input.Title, input.Comment, input.Repeat, input.ID)

	result, err := db.Exec(query, input.Date, input.Title, input.Comment, input.Repeat, input.ID)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v\n", err)
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при проверке обновления задачи: %v\n", err)
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при проверке обновления задачи")
		return
	}
	if rowsAffected == 0 {
		log.Println("Ошибка: Задача не найдена")
		utils.WriteError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	updatedTask := utils.Task{
		ID:      input.ID,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	log.Printf("Возвращаемый JSON: %+v\n", updatedTask)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		log.Printf("Ошибка при кодировании JSON: %v\n", err)
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при отправке данных")
		return
	}

}
