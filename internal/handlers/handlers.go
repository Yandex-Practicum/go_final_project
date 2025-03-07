package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_final-project/internal/db"
	"go_final-project/internal/task"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetTasksHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, `{"error":"Only GET method is supported"}`, http.StatusMethodNotAllowed)
			return
		}
		search := req.URL.Query().Get("search")
		var dateSearch string

		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			dateSearch = parsedDate.Format("20060102")
		}

		tasks, err := db.GetTasks(dbase, search, dateSearch)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []task.Task{}
		}

		tasksList := make([]map[string]string, len(tasks))
		for i, t := range tasks {
			tasksList[i] = map[string]string{
				"id":      strconv.FormatInt(t.ID, 10),
				"date":    t.Date,
				"title":   t.Title,
				"comment": t.Comment,
				"repeat":  t.Repeat,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasksList})
	}
}

func AddTaskHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}
		var newTask task.Task
		err := json.NewDecoder(req.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, `{"error":"Invalid JSON format"}`, http.StatusBadRequest)
			return
		}

		if newTask.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		now := time.Now()
		today := now.Format("20060102")

		if newTask.Date == "" {
			newTask.Date = today
		} else {
			_, err := time.Parse("20060102", newTask.Date)
			if err != nil {
				http.Error(w, `{"error":"Invalid date format, use YYYYMMDD"}`, http.StatusBadRequest)
				return
			}
		}

		taskDate, _ := time.Parse("20060102", newTask.Date)

		if taskDate.Before(now) {
			if newTask.Repeat == "" {
				newTask.Date = today
			} else {
				nextDate, err := NextDate(now, newTask.Date, newTask.Repeat)
				if err != nil || nextDate == "" {
					http.Error(w, `{"error":"Invalid repeat format or no valid next date found"}`, http.StatusBadRequest)
					return
				}
				newTask.Date = nextDate
			}
		}

		id, err := db.AddTask(dbase, &newTask)
		if err != nil {
			http.Error(w, `{"error":"Failed to save task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"id": id})
	}
}

// NextDateHandler determines the next date according to the request parameters.
func NextDateHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		now, err := parseDate(req.URL.Query().Get("now"), time.Now())
		if err != nil {
			http.Error(w, "Invalid 'now' date format. Use YYYYMMDD.", http.StatusBadRequest)
			return
		}

		dateStr := req.URL.Query().Get("date")
		if _, err := parseDate(dateStr, time.Time{}); err != nil {
			http.Error(w, "Invalid 'date' format. Use YYYYMMDD.", http.StatusBadRequest)
			return
		}

		nextDate, err := NextDate(now, dateStr, req.URL.Query().Get("repeat"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, nextDate)
	}
}

func parseDate(dateStr string, defaultDate time.Time) (time.Time, error) {
	if dateStr == "" {
		return defaultDate, nil
	}
	return time.Parse("20060102", dateStr)
}

func jsonResponse(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

// NextDate determines the next date according to the passed parameters.
func NextDate(now time.Time, dateStr string, repeatStr string) (string, error) {
	date, err := parseDate(dateStr, time.Time{})
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	if repeatStr == "" {
		return "", fmt.Errorf("the repeat parameter is not specified")
	}

	firstRepeatLetter := string(repeatStr[0])

	switch firstRepeatLetter {
	case "d":
		return dailyRepeat(now, date, repeatStr)
	case "y":
		return yearlyRepeat(now, date)
	case "w":
		return weeklyRepeat(now, date, repeatStr)
	case "m":
		return monthlyRepeat(now, date, repeatStr)
	default:
		return "", fmt.Errorf("unsupported repeat parameter: %s", repeatStr)
	}
}

// dailyRepeat daily repetition processing 'd'.
func dailyRepeat(now, date time.Time, repeatStr string) (string, error) {
	parts := strings.Fields(repeatStr)
	if len(parts) < 2 {
		return "", fmt.Errorf("the interval in days is not specified")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("error converting to a number: %w", err)
	}
	if days < 1 || days > 400 {
		return "", fmt.Errorf("the minimum and maximum allowed number of days is 1 and 400")
	}

	var nextDate time.Time
	nextDate = date.AddDate(0, 0, days)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}
	return nextDate.Format("20060102"), nil
}

// yearlyRepeat processing the annual recurrence 'y'.
func yearlyRepeat(now, date time.Time) (string, error) {
	var nextDate time.Time
	nextDate = date.AddDate(1, 0, 0)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format("20060102"), nil
}

// weeklyRepeat processing weekly repetition 'w'.
func weeklyRepeat(now, date time.Time, repeatStr string) (string, error) {
	parts := strings.Fields(repeatStr)
	if len(parts) < 2 {
		return "", fmt.Errorf("the interval in weeks is not specified")
	}
	var daysOfWeek []int
	daysList := strings.Split(parts[1], ",")
	for _, d := range daysList {
		dayInt, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil {
			return "", fmt.Errorf("error converting to a number: %w", err)
		}
		if dayInt < 1 || dayInt > 7 {
			return "", fmt.Errorf("day of week must be between 1 and 7")
		}
		daysOfWeek = append(daysOfWeek, dayInt)
	}

	var nextDate time.Time
	nextDate = date
	if nextDate.Before(now) {
		nextDate = now.AddDate(0, 0, 1)
	}
	for {
		weekday := int(nextDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for _, day := range daysOfWeek {
			if day == weekday {
				return nextDate.Format("20060102"), nil
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

// monthlyRepeat monthly repeat processing 'm'.
func monthlyRepeat(now, date time.Time, repeatStr string) (string, error) {
	parts := strings.Fields(repeatStr[2:])
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid repeat parameter format for 'm'")
	}

	days, months, err := parseMonthlyRepeat(parts)
	if err != nil {
		return "", fmt.Errorf("parseMonthlyRepeat: %w", err)
	}

	var nextDate time.Time
	year := now.Year()

	if len(parts) == 2 {
		for {
			month := nextDate.Month()
			if len(months) > 0 && !contains(months, int(month)) {
				nextDate = nextDate.AddDate(0, 1, 0)
				continue
			}
			daysInMonth := lastDayOfMonth(nextDate).Day()
			for _, day := range days {
				if day == -1 {
					day = daysInMonth
				} else if day == -2 {
					day = daysInMonth - 1
				} else if day < 1 || day > daysInMonth {
					continue
				}

				candidate := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
				if candidate.After(now) {
					return candidate.Format("20060102"), nil
				}
			}
			nextDate = nextDate.AddDate(0, 1, 0)
		}
	} else {
		var validDays []int
		for _, dayInt := range days {
			validDays = append(validDays, dayInt)
		}
		nextDate = date
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if nextDate.After(now) {
				for _, day := range validDays {
					if (day == -1 && nextDate.Day() == lastDayOfMonth(nextDate).Day()) ||
						(day == -2 && nextDate.Day() == lastDayOfMonth(nextDate).AddDate(0, 0, -1).Day()) ||
						(day > 0 && nextDate.Day() == day) {
						return nextDate.Format("20060102"), nil
					}
				}
			}
		}
	}
}

// parseMonthlyRepeat parses days and months for monthly repetition.
func parseMonthlyRepeat(parts []string) ([]int, []int, error) {
	var days []int
	var months []int

	for _, p := range strings.Split(parts[0], ",") {
		day, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || (day < -2 || day > 31) || day == 0 {
			return nil, nil, fmt.Errorf("invalid day format: %s", p)
		}
		days = append(days, day)
	}

	if len(parts) > 1 {
		for _, p := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil || (month < 1 || month > 12) {
				return nil, nil, fmt.Errorf("invalid month format: %s", p)
			}
			months = append(months, month)
		}
	}
	return days, months, nil
}

// contains checks whether an element is contained in an array.
func contains(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// lastDayOfMonth returns the last day of the month.
func lastDayOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location())
}
