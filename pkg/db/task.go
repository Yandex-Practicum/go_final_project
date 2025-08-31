package db

import "fmt"

type Task struct {
	ID      int64  `db:"id" json:"id,omitempty"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	q := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB().Exec(q, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	const q = `SELECT id, date, title, comment, repeat
	           FROM scheduler
	           ORDER BY date ASC, id ASC
	           LIMIT ?`

	rows, err := DB().Query(q, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	list := make([]*Task, 0, limit)
	for rows.Next() {
		var t Task
		if err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func GetTask(id string) (*Task, error) {
	row := DB().QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	)
	var t Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(t *Task) error {
	res, err := DB().Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID,
	)
	if err != nil {
		return err
	}
	var n int64
	n, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB().Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB().Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
