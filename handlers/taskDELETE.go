package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func TaskDELETE(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	fmt.Println(id)
	newID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Printf("Error: не парсится ID:%v", id)
		return
	}

	err = dbHelper.DeleteTask(newID)
	if err != nil {
		http.Error(w, `{"error": "ошибка удаления"}`, http.StatusBadRequest)
		log.Printf("Error: ошибка удаления")
		return
	}

	resp := map[string]interface{}{}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
