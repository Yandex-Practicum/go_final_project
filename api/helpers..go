package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Masteker/go_final_project/database/db"
)

// helpers.go содержит вспомогательные функции для работы других хендлеров

// Общие перменные для пакета
var (
	dbh        db.DBHandler
	dateFormat string
)

// ApiInit инициплизирует переменные используемые в пакете api, зависящие от переменных среды и других пакетов
func ApiInit() {
	dbh = db.DBHandl
	dateFormat = os.Getenv("TODO_DATEFORMAT")
}

// isID возвращает true если переданная строка содержит только символы, которые могут находится в строке ID в базе данных.
func isID(id string) bool {
	isID, _ := regexp.Match("[0-9]+", []byte(id))
	return isID
}

// writeErr пишет ошибку в response в формате JSON и статус запроса BadRequest
func writeErr(err error, w http.ResponseWriter) {
	log.Println(err)
	errResp := map[string]string{
		"error": err.Error(),
	}
	resp, err := json.Marshal(errResp)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// writeEmptyJson пишет в response пустой JSON {} и статус запроса OK
func writeEmptyJson(w http.ResponseWriter) {
	okResp := map[string]string{}
	resp, err := json.Marshal(okResp)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
