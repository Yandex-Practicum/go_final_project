package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Обработчик для API запроса
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Преобразуем строку "now" в тип time
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректная дата 'now'", http.StatusBadRequest)
		return
	}

	// Вызываем функцию NextDate для вычисления следующей даты
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем результат в ответ
	if next != "" {
		w.Write([]byte(next))
	} else {
		w.Write([]byte("Задача выполнена"))
	}
}

// NextDate вычисляет следующую дату для задачи в соответствии с указанным правилом повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим исходную дату из формата "20060102"
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %v", err)
	}

	// Проверяем правило повторения
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Если правило повторения - ежегодное "y"
	if repeat == "y" {
		// Добавляем 1 год
		nextDate := taskDate.AddDate(1, 0, 0)

		// Проверяем, что полученная дата больше текущей
		for !nextDate.After(now) {
			// Если дата еще не после текущей, добавляем еще 1 год
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		// Возвращаем дату в нужном формате
		return nextDate.Format("20060102"), nil
	}

	// Если правило повторения - перенос на несколько дней "d <число>"
	if strings.HasPrefix(repeat, "d") {
		// Разделяем строку на команду и количество дней
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", errors.New("неправильный формат правила 'd <число>'")
		}
		var days int
		_, err := fmt.Sscanf(parts[1], "%d", &days)
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверное количество дней или превышен лимит")
		}

		// Переносим задачу на указанное количество дней
		nextDate := taskDate.AddDate(0, 0, days)

		// Проверяем, что полученная дата больше или равна текущей
		if !nextDate.After(now) {
			// Если дата не больше текущей, продолжаем добавлять дни
			for !nextDate.After(now) {
				// Пересчитываем следующую дату, прибавляя дни
				nextDate = nextDate.AddDate(0, 0, days)
			}
		}

		// Возвращаем дату в нужном формате
		return nextDate.Format("20060102"), nil
	}

	// Если правило не поддерживается
	return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
}
