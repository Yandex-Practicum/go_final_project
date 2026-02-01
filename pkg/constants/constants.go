package constants

const (
	Schema string = `
			CREATE TABLE IF NOT EXISTS scheduler (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			date CHAR(8) NOT NULL DEFAULT "",
    			title VARCHAR(255) NOT NULL DEFAULT "",
    			comment TEXT,
    			repeat VARCHAR(128) NOT NULL DEFAULT ""
			);

			CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`
	QueryAddTask                          = `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4);`
	QueryGetTaskListWithTaskFilter string = `SELECT * FROM scheduler WHERE title LIKE ($1) OR comment LIKE ($1) ORDER BY date ($2) LIMIT ($3) OFFSET ($4);`
	QueryGetTaskListWithDateFilter string = `SELECT * FROM scheduler WHERE date = ($1) LIMIT ($2) OFFSET ($3);`
	QueryGetTaskById               string = `SELECT * FROM scheduler WHERE id = $1;`
	QueryUpdateTask                       = `UPDATE scheduler 
							   SET
 								  date = COALESCE(NULLIF($2, ''), date),
 								  title = COALESCE(NULLIF($3, ''), title), 
 								  comment = COALESCE($4, comment),
 								  repeat = COALESCE(NULLIF($5, ''), repeat)
 							   WHERE id = $1`
	QueryDeleteTask string = `DELETE FROM scheduler WHERE id = $1`
)

const (
	DateFormat    string = "20060102"
	DaySign       string = "d"
	MonthSign     string = "m"
	WeekSign      string = "w"
	YearSign      string = "y"
	DaysMinValue  int    = 1
	DaysMaxValue  int    = 400
	WeeksMinValue int    = 1
	WeeksMaxValue int    = 7
	SortASC       string = "ASC"
	SortDESC      string = "DESC"
	DefaultLimit  int64  = 10
)
