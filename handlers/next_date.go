package handlers

import (
	"fmt"
	"net/http"

	"go_final_project/date"
)

// GetNextDateHandler - возвращает значение новой даты, если оно валидно.
func GetNextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	dateToTask := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверяем корректность формата входящего времени
	nowTime, err := date.Validation(now)
	if err != nil {
		fmt.Println("Ошибка конвертации входящего времени nowTime. ", err)
	} else {
		// Получение следующей даты
		res, err := date.NextDate(nowTime, dateToTask, repeat)
		if err != nil {
			fmt.Println("Ошибка получения NextDate. ", err)
		}

		_, err = w.Write([]byte(res))
		if err != nil {
			fmt.Println("Ошибка формирования ответа. ", err)
		}
	}
}
