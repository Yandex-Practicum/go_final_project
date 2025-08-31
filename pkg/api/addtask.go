package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		delTaskHandler(w, r)
	default:
		writeJSON(w, map[string]any{"error": "method not allowed"})
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	t.Title = strings.TrimSpace(t.Title)
	t.Comment = strings.TrimSpace(t.Comment)
	t.Repeat = strings.TrimSpace(t.Repeat)
	if t.Title == "" {
		writeJSON(w, map[string]any{"error": "title is required"})
		return
	}
	if err := fixDate(&t); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{"id": id})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]any{"error": "missing id"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]any{"error": "task not found"})
		return
	}
	writeJSON(w, map[string]string{
		"id":      strconv.FormatInt(t.ID, 10),
		"date":    t.Date,
		"title":   t.Title,
		"comment": t.Comment,
		"repeat":  t.Repeat,
	})
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	idVal, ok := body["id"]
	if !ok {
		writeJSON(w, map[string]any{"error": "missing id"})
		return
	}
	id, err := parseID(idVal)
	if err != nil {
		writeJSON(w, map[string]any{"error": "bad id"})
		return
	}
	t := db.Task{
		ID:      id,
		Date:    toStr(body["date"]),
		Title:   strings.TrimSpace(toStr(body["title"])),
		Comment: strings.TrimSpace(toStr(body["comment"])),
		Repeat:  strings.TrimSpace(toStr(body["repeat"])),
	}
	if t.Title == "" {
		writeJSON(w, map[string]any{"error": "title is required"})
		return
	}
	if err := fixDate(&t); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}

func delTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]any{"error": "missing id"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]any{"error": "missing id"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]any{"error": "task not found"})
		return
	}
	if strings.TrimSpace(t.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]any{})
		return
	}
	now := today(time.Now())
	next, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		writeJSON(w, map[string]any{"error": fmt.Sprintf("bad repeat: %v", err)})
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}

func fixDate(t *db.Task) error {
	now := today(time.Now())

	if t.Date == "" {
		t.Date = now.Format(DateFmt)
	} else {
		if _, err := time.Parse(DateFmt, t.Date); err != nil {
			return errors.New("bad date format")
		}
	}

	dt, _ := time.Parse(DateFmt, t.Date)

	if t.Repeat != "" {
		next, err := NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return fmt.Errorf("bad repeat: %v", err)
		}
		if afterDay(now, dt) {
			t.Date = next
		}
	} else {
		if afterDay(now, dt) {
			t.Date = now.Format(DateFmt)
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(v)
}

func parseID(v any) (int64, error) {
	switch x := v.(type) {
	case float64:
		return int64(x), nil
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(x), 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, errors.New("bad id")
	}
}

func toStr(v any) string {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", x)
	}
}
