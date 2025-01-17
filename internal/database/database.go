package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const TaskLimitShow = 20 // Предел отображаемых задач

type Scheduler struct {
	Id      string `db:"id" json:"id"`           // автоинкрементный идентификатор
	Date    string `db:"date" json:"date"`       // дата в формате YYYYMMDD
	Title   string `db:"title" json:"title"`     // заголовок задачи
	Comment string `db:"comment" json:"comment"` // комментарий к задаче
	Repeat  string `db:"repeat" json:"repeat"`   // правила повторений (максимум 128 символов)
}

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func New() (*sql.DB, error) {
	// Проверим наличие файла с базой
	dbFile, install, err := Check()
	if err != nil {
		return nil, err
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		fmt.Println("База данных не найдена. Создаем...")
		CreateDB(dbFile)
	}
	return sql.Open("sqlite3", "scheduler.db")
}

func (s TaskStore) Add(task Scheduler) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, errors.New("ошибка добавления задачи в базу")
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("не получил идентификатор последней добавленной записи")
	}
	return int(id), nil
}

func (s TaskStore) GetTasks() ([]Scheduler, error) {
	// TODO Добавьте возможность выбрать задачи через строку поиска

	// SQL-запрос для получения ближайших задач
	query := `SELECT * FROM scheduler ORDER BY date LIMIT :limit`
	rows, err := s.db.Query(query, sql.Named("limit", TaskLimitShow))
	if err != nil {
		return []Scheduler{}, errors.New("ошибка в запросе к базе данных")
	}
	defer rows.Close()

	// Массив для хранения задач
	var tasks []Scheduler

	// Чтение данных из базы данных и создание массива задач
	for rows.Next() {
		var task Scheduler
		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []Scheduler{}, errors.New("ошибка при чтении (Rows.Scan) базы данных")
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return []Scheduler{}, err
	}

	return tasks, nil
}

func (s TaskStore) GetTask(id string) (Scheduler, error) {

	// Выполнение SQL-запроса для получения задачи по Id
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
	row := s.db.QueryRow(query, sql.Named("id", id))

	// Чтение данных из базы данных и создание объекта задачи
	var task Scheduler
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Scheduler{}, errors.New("ошибка при чтении (Row.Scan) базы данных")
	}
	return task, nil
}

func (s TaskStore) Update(task Scheduler) error {
	// обновляем данные задачи в базе
	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	_, err := s.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id),
	)
	if err != nil {
		return errors.New("ошибка обновления задачи в базе")
	}

	return nil
}

func (s TaskStore) Delete(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"
	_, err := s.db.Exec(query, sql.Named("id", id))
	return err
}

// TaskExists Функция для проверки существования записи по Id в таблице scheduler
func (s TaskStore) TaskExists(id string) (bool, error) {
	// SQL-запрос для проверки существования записи
	query := `SELECT COUNT(*) FROM scheduler WHERE id = ?`

	// Выполняем запрос
	var count int
	err := s.db.QueryRow(query, id).Scan(&count)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return count > 0, nil
	}
}

func Check() (string, bool, error) {
	appPath, err := os.Executable()
	if err != nil {
		return "", false, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	if err != nil {
		return dbFile, true, nil
	}
	return dbFile, false, nil
	//TODO Реализуйте возможность определять путь к файлу базы данных через переменную окружения. Для этого сервер должен получать значение переменной окружения TODO_DBFILE и использовать его в качестве пути к базе данных, если это не пустая строка.
}

func CreateDB(path string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем базу данных и индекс
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX idx_date ON scheduler(date);
	`); err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных успешно создана и настроена.")
	// Закрываем соединение с базой данных
	defer db.Close()
}
