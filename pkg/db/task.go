package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Task описывает задачу
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`    // формат YYYYMM02
	Title   string `json:"title"`   // обязательное
	Comment string `json:"comment"` // опционально
	Repeat  string `json:"repeat"`  // строка правила
}

// AddTask вставляет задачу в таблицу scheduler и возвращает её ID
func AddTask(task *Task) (int64, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (?, ?, ?, ?)
    `
	res, err := DB.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("db insert error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("db last insert id error: %w", err)
	}
	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	baseQuery := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
    `
	var args []interface{}
	where := ""
	// Попытка распознать DD.MM.YYYY
	if len(search) == len("02.01.2006") && strings.Count(search, ".") == 2 {
		// переводим в YYYYMMDD
		parts := strings.Split(search, ".")
		date := parts[2] + parts[1] + parts[0]
		where = " WHERE date = ?"
		args = append(args, date)
	} else if search != "" {
		where = " WHERE title LIKE ? OR comment LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	orderLimit := " ORDER BY date ASC LIMIT ?"
	args = append(args, limit)

	query := baseQuery + where + orderLimit

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.Tasks query error: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("db.Tasks scan error: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db.Tasks rows error: %w", err)
	}

	// чтобы JSON с empty slice был [], а не null
	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	// Сначала конвертируем строковый ID в число
	var idInt int64
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// SELECT по PK
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE id = ?
        LIMIT 1
    `
	row := DB.QueryRow(query, idInt)

	var t Task
	var idOut int64
	if err := row.Scan(&idOut, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("db scan error: %w", err)
	}
	// конвертируем обратно в строку
	t.ID = fmt.Sprint(idOut)
	return &t, nil
}

func UpdateTask(t *Task) error {
	// Конвертируем ID
	var idInt int64
	if _, err := fmt.Sscan(t.ID, &idInt); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	// Обновляем все колонки, кроме PK
	query := `
        UPDATE scheduler
        SET date    = ?,
            title   = ?,
            comment = ?,
            repeat  = ?
        WHERE id = ?
    `
	res, err := DB.Exec(query,
		t.Date,
		t.Title,
		t.Comment,
		t.Repeat,
		idInt,
	)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("no task with id %s", t.ID)
	}
	return nil
}
