package db

const (
	CreateTableQuery = `
		CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL ,
			title TEXT NOT NULL,
			comment TEXT NOT NULL,
			"repeat" TEXT NOT NULL CHECK (LENGTH("repeat") <= 128)
		);
		`
	CreateIndex = `
	CREATE INDEX scheduler_date 
	    ON scheduler(date);
`
	InsertData = `
		INSERT INTO scheduler( date, title, comment, "repeat") VALUES (?,?,?,?);
		`

	UpdateData = `
		UPDATE scheduler SET date = ? WHERE id = ?;
		`

	GetTasks = `
	SELECT * FROM scheduler ORDER BY date LIMIT ?;
`

	GetTasksByDate = `
	SELECT * FROM scheduler WHERE date= 	:date ORDER BY date LIMIT :limit;
`

	GetTasksByWords = `
	SELECT * FROM scheduler WHERE LOWER(title) LIKE LOWER(:word) OR LOWER(comment) LIKE LOWER(:word) ORDER BY date LIMIT :limit;
`

	GetTaskByID = `
	SELECT * FROM scheduler WHERE id = :id;
`

	UpdateTask = `
	UPDATE scheduler SET id= :id, date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id;
`

	DeleteTask = `
	DELETE FROM scheduler WHERE id= :id;
`
)
