package repository

const (
	SQLCreateScheduler = `
	CREATE TABLE scheduler (
	    id      INTEGER PRIMARY KEY, 
	    date    CHAR(8) NOT NULL DEFAULT "", 
	    title   TEXT NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat  VARCHAR(128) NOT NULL DEFAULT "" 
	);
	`
	SQLCreateSchedulerIndex = `
	CREATE INDEX scheduler_date_index ON scheduler (date)
	`

	SQLCreateTask = `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`

	SQLGetTasks = `SELECT * FROM scheduler WHERE date >= $1`

	SQLGetTasksBySearchString = `SELECT * FROM scheduler WHERE title LIKE '%$1%' OR comment LIKE '%$1%' ORDER BY date`

	SQLGetTasksByDate = `SELECT * FROM scheduler WHERE date = $1`

	SQLGetTaskById = `SELECT * FROM scheduler WHERE id = $1`

	SQLUpdateTask = `UPDATE scheduler SET date = $2, title = $3, comment = $4, repeat = $5 WHERE id = $1`

	SQLMakeTaskDone = `UPDATE scheduler SET date = $2 WHERE id = $1`

	SQLDeleteTask = `DELETE FROM scheduler WHERE id = $1`
)
