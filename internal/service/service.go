package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ASHmanR17/go_final_project/internal/database"
)

const (
	DateLayout = "20060102"
)

type TaskService struct {
	store database.TaskStore
}

func NewTaskService(store database.TaskStore) *TaskService {
	return &TaskService{store: store}
}

// TaskFromJson считывает JSON из body и десериализует его в структуру task
func (t TaskService) TaskFromJson(body io.Reader) (database.Scheduler, error) {
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(body)
	if err != nil {
		return database.Scheduler{}, errors.New("Ошибка чтения JSON")
	}

	var task database.Scheduler
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		return database.Scheduler{}, errors.New("Ошибка декодирования JSON")
	}
	return task, nil
}

// CheckTask проверяет правильность полей задачи
func (t TaskService) CheckTask(task database.Scheduler) (database.Scheduler, error) {
	// Получаем текущую дату и время
	currentDate := time.Now()
	// Проверим наличие заголовка Title
	if task.Title == "" {
		return database.Scheduler{}, errors.New("задача не имеет заголовка")
	}
	if task.Date == "" { // Проверка пустой даты. Если содержит пустую строку, берётся сегодняшнее число
		// Форматируем текушую дату в формате "YYYYMMDD" и присваиваем её полю Date
		task.Date = currentDate.Format(DateLayout)
	}
	// Проверим корректность даты
	// Регулярное выражение для проверки формата даты
	validFormat := regexp.MustCompile(`^\d{8}$`)
	// Проверка формата даты
	taskTime, err := time.Parse(DateLayout, task.Date)
	if err != nil || !validFormat.MatchString(task.Date) {
		return database.Scheduler{}, errors.New("некорректная дата")
	}

	// вычисляем следующую дату и проверим правила повторения
	nextDate, err := NextDate(currentDate, task.Date, task.Repeat)
	if err != nil {
		return database.Scheduler{}, err
	}

	//если дата задачи меньше сегодняшней, подставляем дату, вычисленную в NextDate
	nowDay := currentDate.Truncate(24 * time.Hour)

	if taskTime.Before(nowDay) {
		task.Date = nextDate
	}

	return task, nil
}

func (t TaskService) AddTask(task database.Scheduler) (int, error) {
	id, err := t.store.Add(task)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// NextDate вычисляет следующую дату для задачи в соответствии с указанным правилом repeat
// Функция возвращает следующую дату в формате DateLayout и ошибку.
// Возвращаемая дата должна быть больше даты, указанной в переменной now.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Преобразуем строку даты в формат time.Time
	startDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	// Обработка правил повторения
	var result time.Time
	var daysToAdd int
	switch {
	case repeat == "": // Если правило не указано
		//Если дата меньше сегодняшнего числа, возвращаем сегодня
		if startDate.Before(now) {
			return now.Format(DateLayout), nil
		}
		// Если правило не указано, возвращаем стартовую date
		return date, nil
	case strings.HasPrefix(repeat, "d "):
		// Преобразуем правило в число дней
		daysToAdd, err = strconv.Atoi(repeat[2:])
		if err != nil || daysToAdd > 400 || daysToAdd < 1 {
			return "", errors.New("некорректное число дней в правиле повторения")
		}
		result = addDays(startDate, daysToAdd, now)

	case repeat == "y":
		result = addYear(startDate, now)
	// TODO добавить правила для недель и месяцев
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
	return result.Format(DateLayout), nil

}

func (t TaskService) GetTasks() ([]database.Scheduler, error) {

	tasks, err := t.store.GetTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskService) GetTask(id string) (database.Scheduler, error) {
	task, err := t.store.GetTask(id)
	if err != nil {
		return database.Scheduler{}, err
	}
	return task, nil
}

func (t TaskService) UpdateTask(task database.Scheduler) error {
	//Проверим корректность Id перед обновлением в базе
	id, err := strconv.Atoi(task.Id)
	if err != nil || id <= 0 {
		return errors.New("неверный ID")
	}
	//Проверим существование записи с таким Id перед обновлением в базе
	exist, err := t.store.TaskExists(task.Id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("записи с таким Id не существует")
	}
	// обновляем данные в базе
	return t.store.Update(task)
}

func (t TaskService) DeleteTask(id string) error {
	return t.store.Delete(id)
}

// addDays вычисляет следующую дату в цикле
func addDays(date time.Time, daysToAdd int, now time.Time) time.Time {
	step := daysToAdd
	for {
		nextDate := date.AddDate(0, 0, daysToAdd)
		if nextDate.After(now) {
			return nextDate
		}
		daysToAdd = daysToAdd + step
	}
}

func addYear(date time.Time, now time.Time) time.Time {
	step := 1
	for {
		nextDate := date.AddDate(step, 0, 0)
		if nextDate.After(now) {
			return nextDate
		}
		step++
	}
}
