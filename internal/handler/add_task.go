package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/askoren1/go_final_project/internal/models"
	nextdate "github.com/askoren1/go_final_project/internal/next_date"
)

const Layout = "20060102"

// обработчик HTTP-запроса для добавления новой задачи в систему
func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { //Декодирование JSON в структуру models.Task
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка правила повторения
	switch t.Repeat {
	case "":
		// Правило "" - корректное, ничего не делаем

	case "y":
		// Правило "y" - корректное, ничего не делаем

	default:
		// Проверяем формат "d <число>"
		match, _ := regexp.MatchString(`^d \d+$`, t.Repeat) // Проверка на соответствие регулярному выражению
		if !match {
			response := map[string]string{"error": "Некорректный формат правила повторения"}
			json.NewEncoder(w).Encode(response)
			return
		}

		daysStr := t.Repeat[2:] // Извлекаем число дней
		days, err := strconv.Atoi(daysStr)
		if err != nil || days > 400 {
			response := map[string]string{"error": "Некорректное значение дней в правиле повторения"}
			json.NewEncoder(w).Encode(response)
			return
		}

	}

	//Проверка заголовка задачи
	if t.Title == "" { // проверка title на пустоту
		response := map[string]string{"error": "Не указан заголовок задачи"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Обработка даты
	// dateToday - сегодня, string
	// nowTime - сегодня, time.Time
	// t.Date - входное значение даты, string
	// date2 - входное значение даты, time.Time
	// dateInTable - дата для передачи в таблицу, string

	dateToday := time.Now().Format(Layout)
	nowTime := time.Now()

	var dateInTable string

	if t.Date == "" || t.Date == "today" { // проверка date на пустоту
		dateInTable = dateToday
	} else { // проверка даты на корректность
		date2, err := time.Parse(Layout, t.Date)
		if err != nil {
			response := map[string]string{"error": "Некорректное значение даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		if nowTime.Truncate(24 * time.Hour).UTC().After(date2) { //замена даты, если дата задачи меньше сегодняшней

			if t.Repeat == "" {
				dateInTable = dateToday
			} else {
				dateInTable, _ = nextdate.NextDate(nowTime, t.Date, t.Repeat)
			}

		} else {
			dateInTable = t.Date
		}
	}

	// Добавление задачи в базу данных
	id, err := h.repo.AddTask(dateInTable, t.Title, t.Comment, t.Repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Формирование ответа
	response := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)

}
