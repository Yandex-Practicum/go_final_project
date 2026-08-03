// package service

// import (
// 	"database/sql"
// 	"fmt"
// )

// func GetTaskID(id string) (*Task, error) {

// 	if id == "" {
// 		return nil, fmt.Errorf("Поле id не может быть пустым.")
// 	}

// 	fmt.Println("программа дошла до сюда0")

// 	db, err := sql.Open("sqlite", PathDbManual())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()
// 	var res Task

// 	fmt.Println("программа дошла до сюда1")

// 	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id ;", sql.Named("id", id))
// 	fmt.Println("программа дошла до сюда2")

// 	err = row.Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)
// 	fmt.Println("программа дошла до сюда3")

// 	if err := row.Err(); err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(&res)

// 	return &res, nil

// 	// byteJson, err := json.Marshal(res)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// }

package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetTaskID(id string) (Task, error, int) {

	if id == "" {
		return Task{}, fmt.Errorf("Поле id не может быть пустым."), http.StatusBadRequest
	}

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		return Task{}, err, http.StatusInternalServerError
	}
	defer db.Close()
	var res Task

	intStr, err := strconv.Atoi(id)

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", intStr)).Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)

	if err != nil {
		log.Println(err)
		return Task{}, err, http.StatusInternalServerError
	}

	return res, nil, http.StatusOK

}
