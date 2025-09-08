package db

import (
	"fmt"
	"strings"
	"time"
)

// Task структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает её ID
func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `
		INSERT INTO scheduler(date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Tasks возвращает список ближайших задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Чтобы JSON возвращал [] вместо null
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// TasksFiltered возвращает список задач с фильтром по search и ограничением
func TasksFiltered(search string, limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	search = strings.TrimSpace(search)
	var query string
	var args []interface{}

	// Проверяем, является ли search датой в формате 02.01.2006
	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		// Если парсинг успешен, фильтруем по дате
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date ASC
			LIMIT ?
		`
		args = append(args, t.Format("20060102"), limit)
	} else if search != "" {
		// Поиск подстроки в title или comment
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date ASC
			LIMIT ?
		`
		likePattern := "%" + search + "%"
		args = append(args, likePattern, likePattern, limit)
	} else {
		// Без фильтра, просто первые limit задач
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT ?
		`
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по id
func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	var t Task
	var taskID int64
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&taskID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}
	t.ID = fmt.Sprintf("%d", taskID)
	return &t, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу по id
func DeleteTask(id string) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(next string, id string) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
