package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go_final_project/date"
	"go_final_project/tasks"

	_ "modernc.org/sqlite"
)

// limit - максимальное число выводимых задач
const limit = 10

// StringError - Строковый тип для специальной ошибки
type StringError struct {
	StrEr string `json:"error,omitempty"`
}

// IDType - Структура для ID
type IDType struct {
	ID int64 `json:"id"`
}

// TasksType - Структура для группы записей
type TasksType struct {
	Tasks []tasks.Task `json:"tasks"`
}

// Storage - Структура для хранилища
type Storage struct {
	DB *sql.DB
}

// NewStorage - онструктор для хранилица
func NewStorage(db *sql.DB) Storage {
	return Storage{DB: db}
}

// TodoStorage - переменная для обращения к БД
var TodoStorage Storage

// GetTask - возвращает информацию об одной задаче по входным данным
func (s *Storage) GetTask(data tasks.Task) ([]byte, error) {

	var result []byte
	var returnData tasks.Task
	var row *sql.Rows
	var err error
	var errRes StringError

	// Формируем запрос в базу
	qeryToDB := `SELECT id, date, title, comment, repeat
					FROM scheduler
				WHERE id = ?;`

	row, err = s.DB.Query(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Чтение из БД не состоялась ", err)
		return result, err
	}
	defer row.Close()

	// Укладываем результаты запроса в структуру
	for row.Next() {
		if err := row.Scan(&returnData.ID, &returnData.Date, &returnData.Title, &returnData.Comment, &returnData.Repeat); err != nil {
			return nil, err
		}
		if err != nil {
			fmt.Println("Не удалось записать корректную дату из БД .", err)
			return result, err
		}
	}

	// Проверка на то, что задача вообще была в базе
	if returnData.ID == "" {
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	// Формируем сообщение с информацией о задаче
	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось сформировать вывод в виде JSON ", err)
		return result, err
	}

	return result, nil
}

// AddTask - записывает в БД данные о внесённой задаче
func (s *Storage) AddTask(data tasks.Task) ([]byte, error) {

	var err error
	var result []byte
	var returnData IDType

	// Формируем запрос в базу
	qeryToDB := `INSERT INTO
					scheduler (date, title, comment, repeat)
				VALUES (?, ?, ?, ?);`

	res, err := s.DB.Exec(qeryToDB, data.Date, data.Title, data.Comment, data.Repeat)
	if err != nil {
		fmt.Println("Запись в БД не состоялась ", err)
		return result, err
	}

	// Возвращаем id последней записи
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		return result, err
	}

	returnData.ID = id

	// Формируем сообщение с информацией о задаче
	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось выдать ID последней записи в виде JSON ", err)
		return result, err
	}

	return result, nil
}

// UpdateTask - изменяет в БД данные о внесённой задаче
func (s *Storage) UpdateTask(data tasks.Task) ([]byte, error) {

	var errRes StringError
	var result []byte
	var err error

	// Формируем запрос в базу
	qeryToDB := `UPDATE
					scheduler SET date = ?, title = ?, comment = ?, repeat = ?
				WHERE id = ? ;`

	res, err := s.DB.Exec(qeryToDB, data.Date, data.Title, data.Comment, data.Repeat, data.ID)
	if err != nil {
		fmt.Println("Запись в БД не состоялась ", err)
		return result, err
	}

	// Возвращаем количество затронутых записей
	num, err := res.RowsAffected()

	if err != nil {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	if num == 0 {
		fmt.Println("ID последней записи в БД не удалось получить, num = 0 ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	// Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// DeleteTask - удаляет из БД данные об одной задаче
func (s *Storage) DeleteTask(data tasks.Task) ([]byte, error) {

	var errRes StringError
	var err error
	var result []byte

	// Формируем запрос в базу
	qeryToDB := `DELETE FROM
					scheduler 
					WHERE id = ?;`

	res, err := s.DB.Exec(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Удаление из БД не состоялась ", err)
		return result, err
	}

	// Возвращаем количество затронутых записей
	id, err := res.LastInsertId()

	if err != nil {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	if id == 0 {
		fmt.Println("ID последней записи в БД не удалось получить, num = 0 ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	// Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// DoneTasks - удаляет из БД данные об одной задаче при её выполнении
func (s *Storage) DoneTasks(data tasks.Task) ([]byte, error) {

	var errRes StringError
	var result []byte
	var returnData tasks.Task
	var row *sql.Rows
	var err error

	now := time.Now()

	// Формируем запрос в базу
	qeryToDB := `SELECT id, date, title, comment, repeat
					FROM scheduler
					WHERE id = ?;`

	row, err = s.DB.Query(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Чтение из БД не состоялась ", err)
		return result, err
	}
	defer row.Close()

	// Укладываем результаты запроса в структуру
	for row.Next() {
		if err := row.Scan(&returnData.ID, &returnData.Date, &returnData.Title, &returnData.Comment, &returnData.Repeat); err != nil {
			return nil, err
		}

		if err != nil {
			fmt.Println("Не удалось записать корректную дату из БД .", err)
			return result, err
		}
	}

	// Есть ли правило для повторения задачи
	switch {
	case returnData.Repeat == "":
		result, err = TodoStorage.DeleteTask(returnData)
		if err != nil {
			fmt.Println("Не удалось удалить задачу из БД .", err)
			return result, err
		}
	default:
		returnData.Date, err = date.NextDate(now, returnData.Date, returnData.Repeat)
		if err != nil {
			fmt.Println("Не удалось получить новую дату для задачи .", err)
			return result, err
		}

		result, err = TodoStorage.UpdateTask(returnData)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
			return result, err
		}
	}

	// Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// FindTasks - возвращает информацию о группе последних задач
func (s *Storage) FindTasks(search string) ([]byte, error) {

	var errRes StringError
	var err error
	var task tasks.Task
	var rows *sql.Rows
	var queryToDB, searchDate string
	var result []byte
	returnData := TasksType{Tasks: make([]tasks.Task, 0, 20)}

	if search != "" && len(search) == 10 {
		searchDate, err = date.Convert(search)
		if err != nil {
			fmt.Println("На входе не дата ", err)
		}
	}

	switch {
	// Ищем всё подряд
	case search == "":
		queryToDB = `SELECT id, date, title, comment, repeat
						FROM scheduler
					ORDER BY date LIMIT ?;`

		rows, err = s.DB.Query(queryToDB, limit)
		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}

	// Ищем по дате
	case searchDate != "":
		queryToDB = `SELECT id, date, title, comment, repeat
							FROM scheduler
							WHERE date = ?
						ORDER BY date LIMIT ?;`

		rows, err = s.DB.Query(queryToDB, searchDate, limit)
		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}

	// Ищем по заголовку или комментарию
	default:

		search = fmt.Sprint("%" + search + "%")

		queryToDB = `SELECT id, date, title, comment, repeat
						FROM scheduler 
						WHERE title LIKE ? OR comment LIKE ?
					ORDER BY date LIMIT ?;`

		rows, err = s.DB.Query(queryToDB, search, search, limit)

		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			fmt.Println("Запись в структуру данных из БД не состоялась ", err)
			return result, err
		}

		if err = rows.Err(); err != nil {
			fmt.Println("Не удалось записать корректные данные из БД .", err)
			return result, err
		}

		returnData.Tasks = append(returnData.Tasks, task)
	}

	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось сформировать вывод в виде JSON ", err)
		errRes.StrEr = fmt.Sprintln(err)
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	return result, nil
}
