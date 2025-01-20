package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lezhnevay/go_final_project_lezhnev/configs"
	"github.com/lezhnevay/go_final_project_lezhnev/internal/tasks"
)

type Store struct {
	db *sql.DB
}

// Открываем/создаем БД
func OpenDataBase() *sql.DB {
	// Определяем директорию приложения и проверяем наличие базы данных
	filePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(filePath), "scheduler.db")

	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	log.Println("Путь к базе данных", dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	// создаем таблицу и индекс
	if install {

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128) NOT NULL
			);`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных была создана ранее")
	}
	return db

}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

// Добавление задачи
func (s *Store) CreateTask(t configs.Task) (string, error) {
	var err error

	if t.Title == "" {
		return "", fmt.Errorf(`{"error":"не указан заголовок задачи"}`)
	}

	// Проверяем наличие даты
	if t.Date == "" {
		t.Date = time.Now().Format(configs.DateFormat)
	}

	_, err = time.Parse(configs.DateFormat, t.Date)
	if err != nil {
		return "", fmt.Errorf(`{"error":"некорректный формат даты"}`)
	}
	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if t.Date < time.Now().Format(configs.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := tasks.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return "", fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(configs.DateFormat)
		}
	}

	// Добавляем задачу в базу
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf(`{"error":"не удалось добавить задачу"}`)
	}

	// Возвращаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"не удалось вернуть id новой задачи"}`)
	}
	return fmt.Sprintf("%d", id), nil
}

// Получение листинга задач по фильтрам
func (s *Store) GetTasks(search string) ([]configs.Task, error) {
	var t configs.Task
	var tasks []configs.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", configs.MaxTasks)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, date.Format(configs.DateFormat), configs.MaxTasks)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, search, search, configs.MaxTasks)
	}
	if err != nil {
		return []configs.Task{}, fmt.Errorf(`{"error":"ошибка запроса"}`)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err = rows.Err(); err != nil {
			return []configs.Task{}, fmt.Errorf(`{"error":"ошибка распознавания данных"}`)
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		tasks = []configs.Task{}
	}

	return tasks, nil
}

// Получение задачи по id
func (s *Store) GetTask(id string) (configs.Task, error) {
	var t configs.Task
	if id == "" {
		return configs.Task{}, fmt.Errorf(`{"error":"не указан id"}`)
	}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return configs.Task{}, fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	return t, nil
}

// Редактирование задачи
func (s *Store) UpdateTask(t configs.Task) error {
	// Проверяем, что id в теле не пустой
	if t.ID == "" {
		return fmt.Errorf(`{"error":"не указан id"}`)
	}
	// Проверяем обязательное поле Title
	if t.Title == "" {
		return fmt.Errorf(`{"error":"не указан заголовок задачи"}`)
	}
	// Проверяем дату
	if t.Date == "" {
		t.Date = time.Now().Format(configs.DateFormat)
	}

	_, err := time.Parse(configs.DateFormat, t.Date)
	if err != nil {
		return fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}

	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if t.Date < time.Now().Format(configs.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := tasks.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {

				return fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(configs.DateFormat)
		}
	}

	// Обновляем задачу в базе
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {

		return fmt.Errorf(`{"error":"задача с таким id не найдена"}`) // вот эта вот история не работает, посчитаем измененные ряды
	}

	// Считаем измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`{"error":"не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(`{"error":"задача с таким id не найдена"}`)
	}
	// Если ошибок не возникло, возвращаем nil
	return nil
}

// Выполнение задачи
func (s *Store) TaskDone(id string) error {
	var t configs.Task

	t, err := s.GetTask(id)
	if err != nil {
		return err
	}
	if t.Repeat == "" {

		err := s.DeleteTask(id)
		if err != nil {
			return err
		}

	} else {
		next, err := tasks.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}
		t.Date = next
		err = s.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}

// Удаление задачи из БД
func (s *Store) DeleteTask(id string) error {
	// Проверяем, что id не пустой
	if id == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(`{"error":"Не удалось удалить задачу"}`)
	}
	// Считаем измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf(`{"error":"Не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {

		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`)
	}
	// Если ошибок нет, возвращаем nil
	return nil
}
