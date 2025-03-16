package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sandrinasava/go_final_project/internal/models"
	"github.com/sandrinasava/go_final_project/internal/services"
	_ "modernc.org/sqlite"
)

func SendErrorResponse(res http.ResponseWriter, message string, statusCode int) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(map[string]string{"error": message})
}

func CheckTaskAndFindDate(task models.Task) (string, error) {

	if task.Title == "" {
		return "", fmt.Errorf("параметр title пустой")
	}

	var date time.Time
	if task.Date != "" {
		date, err := time.Parse(services.Format, task.Date)
		if err != nil {
			return "", fmt.Errorf("указан неверный формат даты %v", err)
		}
		if date.Format(services.Format) < time.Now().Format(services.Format) {

			if task.Repeat != "" {

				now := time.Now().Format(services.Format)
				dateStr, err := services.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return "", fmt.Errorf("ошибка поиска NextDate %v", err)
				}

				log.Printf("следующая дата = %s", dateStr)
				return dateStr, nil

			} else {
				date = time.Now()
				return date.Format(services.Format), nil

			}
		}
		return date.Format(services.Format), nil
	}

	date = time.Now()
	return date.Format(services.Format), nil

}
