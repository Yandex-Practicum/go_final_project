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

	SQLGetTaskById = `SELECT * FROM scheduler WHERE id = $1`

	SQLUpdateTask = ``
)
