package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSON записывает данные в ответ HTTP в формате JSON.
// Параметры:
// w - объект http.ResponseWriter для отправки ответа клиенту;
// status - HTTP-статус-код, который будет отправлен в ответе;
// data - произвольные данные, которые нужно закодировать в JSON и отправить.
// Возвращает:
// ошибку, если кодирование в JSON или запись в ResponseWriter не удались, nil в случае успешного выполнения.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Устанавливаем заголовки и статус заранее
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Обрабатываем nil-данные: отправляем пустой объект {}
	if data == nil {
		_, err := w.Write([]byte("{}"))
		return err
	}

	// Создаём энкодер с экранированием HTML-символов (безопасность)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)

	// Кодируем данные в JSON
	err := encoder.Encode(data)
	if err != nil {
		log.Printf("JSON encoding error: %v", err)
		return err
	}

	return nil
}
