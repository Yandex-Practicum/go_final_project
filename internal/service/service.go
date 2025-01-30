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
		return "", errors.New("некорректная1 дата")
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

	case strings.HasPrefix(repeat, "m "):
		// Разделяем правило на два массива
		var monthsPart string = "1,2,3,4,5,6,7,8,9,10,11,12"
		ruleParts := strings.Split(repeat[2:], " ")
		if len(ruleParts) == 2 {
			monthsPart = ruleParts[1]
			//return "", errors.New("некорректное правило повторения")
		}
		daysPart := ruleParts[0]

		// Переносим дату на ближайший указанный день месяца
		result, err = findNearestDateWithMonths(daysPart, monthsPart, now)
		if err != nil {
			return "", err
		}

	case repeat == "y":
		result = addYear(startDate, now)
	// TODO добавить правила для недель и месяцев
	case strings.HasPrefix(repeat, "w "):
		// Преобразуем правило в дни недели
		var day int
		var days []string
		currentDay := now.Day()

		closestDay := -1
		days = strings.Split(repeat[2:], " ")
		// Находим ближайший день недели из массива
		for _, d := range days {
			day, err = strconv.Atoi(d)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("некорректный день недели в правиле повторения")
			}
			if day > currentDay {
				if closestDay < 0 || day < closestDay {
					closestDay = day
				}
			}
		}
		if closestDay < 0 {
			closestDay = 1
		}
		// Переносим дату на ближайший указанный день недели
		result = findNearestDay(startDate, closestDay)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
	return result.Format(DateLayout), nil
}

func (t TaskService) GetTasks(search string) ([]database.Scheduler, error) {
	// Если строка поиска пуста, возвращаем все задачи
	if search == "" {
		tasks, err := t.store.GetTasks()
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}
	date, err := time.Parse("02.01.2006", search)
	if err == nil {
		// Если search является датой, получаем задачи на эту дату
		tasks, err := t.store.GetTasksByDate(date.Format(DateLayout))
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}
	tasks, err := t.store.GetTasksBySearch(search)
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

	// обновляем данные в базе
	return t.store.Update(task)
}

func (t TaskService) DeleteTask(id string) error {
	return t.store.Delete(id)
}

func (t TaskService) DoneTask(id string) error {
	// Получим из базы задачу по Id
	task, err := t.GetTask(id)
	if err != nil {
		return err
	}

	// Если правила повторения нет, удаляем задачу из базы и заканчиваем функцию
	if task.Repeat == "" {
		err := t.DeleteTask(task.Id)
		if err != nil {
			return err
		}
		return nil
	}
	// Получаем текущую дату и время
	currentDate := time.Now()
	// вычисляем следующую дату и заодно проверим правила повторения
	nextDate, err := NextDate(currentDate, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	// обновим в базе задачу с новой датой
	task.Date = nextDate
	err = t.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
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

// findNearestDay, которая определяет дату, ближайшую к startDate,
// и которая приходится на день недели closestDay, где 1 — это понедельник, а 7 — воскресенье:
func findNearestDay(startDate time.Time, closestDay int) time.Time {
	// Получаем текущий день недели (1 - понедельник, 7 - воскресенье)
	currentDay := int(startDate.Weekday())

	// Если текущий день недели равен 0 (воскресенье), то присваиваем ему значение 7
	if currentDay == 0 {
		currentDay = 7
	}

	// Вычисляем разницу между текущим днем недели и ближайшим днем
	dayDifference := (closestDay - currentDay + 7) % 7

	// Если разница равна 0, то ближайший день - это текущий день
	if dayDifference == 0 {
		return startDate
	}

	// В противном случае, добавляем разницу к текущей дате
	nearestDay := startDate.AddDate(0, 0, dayDifference)

	return nearestDay
}

func lastDayOfMonth(monthNum int, now time.Time) int {

	date := time.Date(now.Year(), time.Month(monthNum+1), 0, 0, 0, 0, 0, time.UTC)
	return date.Day()
}
func secondLastDayOfMonth(monthNum int, now time.Time) int {

	date := time.Date(now.Year(), time.Month(monthNum+1), 0, 0, 0, 0, 0, time.UTC)
	secondLastDay := date.AddDate(0, 0, -1)
	return secondLastDay.Day()
}

func getDate(day int, month int, now time.Time) (time.Time, error) {

	// Создаем дату
	date := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)

	//// Проверяем, что дата корректна
	//if date.Month() != time.Month(month) || date.Day() != day {
	//	return time.Time{}, errors.New("некорректная дата")
	//}

	return date, nil
}

func findNearestDateWithMonths(days string, months string, now time.Time) (time.Time, error) {
	// Разделяем строки на массивы
	daysArray := strings.Split(days, ",")
	monthsArray := strings.Split(months, ",")

	// Создаем список дат
	var dates []time.Time
	for _, day := range daysArray {
		dayNum, err := strconv.Atoi(day)
		if err != nil || dayNum < -2 || dayNum > 31 || dayNum == 0 {
			return time.Time{}, errors.New("некорректный день месяца в правиле повторения")
		}
		for _, month := range monthsArray {
			monthNum, err := strconv.Atoi(month)
			if err != nil || monthNum < 1 || monthNum > 12 {
				return time.Time{}, errors.New("некорректный месяц в правиле повторения")
			}

			if dayNum == -1 {
				// Получаем последний день месяца
				dayNum = lastDayOfMonth(monthNum, now)
			}
			if dayNum == -2 {
				dayNum = secondLastDayOfMonth(monthNum, now)
			}

			date, err := getDate(dayNum, monthNum, now)
			if err != nil {
				return time.Time{}, err
			}
			dates = append(dates, date)
		}
	}

	// Находим ближайшую дату
	nearestDate := time.Time{}
	for _, date := range dates {
		if date.After(now) {
			if nearestDate.IsZero() {
				nearestDate = date
				continue
			}
		}
		if date.After(now) && date.Before(nearestDate) {
			nearestDate = date
		}
	}

	return nearestDate, nil
}
