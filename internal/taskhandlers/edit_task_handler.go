package taskhandlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"go_final_project/internal/domain/entities"
)

// EditTaskHandler - handler for updating a task (PUT /api/task)
func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
			return
		}

		// Decode JSON into a map[string]string to handle ID as a string
		var rawData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			http.Error(w, `{"error":"Error parsing JSON"}`, http.StatusBadRequest)
			return
		}

		// Convert ID to int64
		id, err := strconv.ParseInt(rawData["id"], 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Invalid ID format"}`, http.StatusBadRequest)
			return
		}

		// Create a JSON object and pass it into `r.Body`
		newBody, _ := json.Marshal(entities.Task{
			ID:      id,
			Date:    rawData["date"],
			Title:   rawData["title"],
			Comment: rawData["comment"],
			Repeat:  rawData["repeat"],
		})

		// Replace `r.Body` with the new JSON (to pass it into `HandleUpdateTask`)
		r.Body = io.NopCloser(bytes.NewReader(newBody))

		// Call `HandleUpdateTask` without passing `task`
		HandleUpdateTask(w, r, db)
	}
}
