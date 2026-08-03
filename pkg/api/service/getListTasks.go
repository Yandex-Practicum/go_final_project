// package service

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"

// 	_ "modernc.org/sqlite"
// )

// func GetTask(limit int, w http.ResponseWriter, r *http.Request) ([]byte, error) {

// 	type TasksResp struct {
// 		Tasks []Task `json:"tasks"`
// 	}

// 	db, err := sql.Open("sqlite", PathDbManual())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	search := r.FormValue("search")
// 	if search == "" {
// 		return nil, nil
// 	} else {
// 		db.Exec("SELECT * FROM scheduler WHERE title, comment LIKE ':search'",
// 			sql.Named("search", search))
// 	}

// 	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit ;",
// 		sql.Named("limit", limit))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var res Task
// 	var resTask TasksResp

// 	for rows.Next() {
// 		err := rows.Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resTask.Tasks = append(resTask.Tasks, res)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	if resTask.Tasks == nil {
// 		resTask.Tasks = []Task{}
// 	}

// 	jsonbyte, err := json.Marshal(resTask)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonbyte, nil

// }

// package service

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	_ "modernc.org/sqlite"
// )

// func GetTask(limit int, w http.ResponseWriter, r *http.Request) ([]byte, error) {

// 	type TasksResp struct {
// 		Tasks []Task `json:"tasks"`
// 	}

// 	db, err := sql.Open("sqlite", PathDbManual())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	search := r.FormValue("search")

// 	var rows *sql.Rows
// 	if search == "" {
// 		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit ;",
// 			sql.Named("limit", limit))
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rows.Close()

// 	} else {
// 		// rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE ':search'",
// 		// rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE ':search';",
// 		rows, err = db.Query(fmt.Sprint("SELECT * FROM scheduler WHERE title LIKE '%", search, "%';"))
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rows.Close()
// 	}

// 	var res Task
// 	var resTask TasksResp

// 	for rows.Next() {
// 		err := rows.Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resTask.Tasks = append(resTask.Tasks, res)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	if resTask.Tasks == nil {
// 		resTask.Tasks = []Task{}
// 	}

// 	jsonbyte, err := json.Marshal(resTask)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonbyte, nil

// }

package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func GetListTask(limit int, w http.ResponseWriter, r *http.Request) ([]byte, error, int) {

	type TasksResp struct {
		Tasks []Task `json:"tasks"`
	}

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer db.Close()

	search := r.FormValue("search")

	var rows *sql.Rows
	if search == "" {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit ;",
			sql.Named("limit", limit))
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		defer rows.Close()
	} else {
		midleLayout := "02.01.2006"
		tims, err := time.Parse(midleLayout, search)
		timesis := tims.Format(LayoutDate)
		if err == nil {
			rows, err = db.Query(fmt.Sprint("SELECT * FROM scheduler WHERE date LIKE '%", timesis, "%';"))
			if err != nil {
				return nil, err, http.StatusInternalServerError
			}
			defer rows.Close()
		} else {
			rows, err = db.Query(fmt.Sprint("SELECT * FROM scheduler WHERE title || comment LIKE '%", search, "%';"))
			if err != nil {
				return nil, err, http.StatusInternalServerError
			}
			defer rows.Close()
		}

	}

	var res Task
	var resTask TasksResp

	for rows.Next() {
		err := rows.Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		resTask.Tasks = append(resTask.Tasks, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if resTask.Tasks == nil {
		resTask.Tasks = []Task{}
	}

	jsonbyte, err := json.Marshal(resTask)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return jsonbyte, nil, http.StatusOK

}
