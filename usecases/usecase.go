package usecases

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlexJudin/go_final_project/repository"
	"github.com/AlexJudin/go_final_project/usecases/model"
)

var _ Task = (*TaskUsecase)(nil)

const (
	year = 1

	sundayEU  = 0
	sundayRus = 7
)

type TaskUsecase struct {
	DB repository.Task
}

func NewTaskUsecase(db repository.Task) *TaskUsecase {
	return &TaskUsecase{DB: db}
}

func (t *TaskUsecase) GetNextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is required")
	}

	dateTask, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	repeatString := strings.Split(repeat, " ")

	switch strings.ToLower(repeatString[0]) {
	case "d":
		if len(repeatString) < 2 {
			return "", fmt.Errorf("repeat should be at least two characters for days")
		}

		days, err := parseValue(repeatString[1])
		if err != nil {
			return "", err
		}
		dateTask = addDateTask(now, dateTask, 0, 0, days)
	case "y":
		dateTask = addDateTask(now, dateTask, year, 0, 0)
	case "w":
		if len(repeatString) < 2 {
			return "", fmt.Errorf("repeat should be at least more two characters for weeks")
		}

		dateTask, err = getDateTaskByWeek(now, dateTask, repeatString[1])
		if err != nil {
			return "", err
		}
	case "m":
		if len(repeatString) < 2 {
			return "", fmt.Errorf("repeat should be at least two characters for month")
		}

		dateTask, err = getDateTaskByMonth(now, dateTask, repeatString)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid character")
	}

	return dateTask.Format("20060102"), nil
}

func (t *TaskUsecase) CreateTask(task *model.Task, pastDay bool) (*model.TaskResp, error) {
	if pastDay {
		nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return nil, err
		}

		task.Date = nextDate
	}

	taskId, err := t.DB.CreateTask(task)
	if err != nil {
		return nil, err
	}

	taskResp := model.NewTaskResp(taskId)

	return taskResp, nil
}

func (t *TaskUsecase) GetTasks(searchString string) (model.TasksResp, error) {
	date, err := time.Parse("02.01.2006", searchString)
	if err == nil {
		return t.DB.GetTasksByDate(date)
	}

	if searchString != "" {
		return t.DB.GetTasksBySearchString(searchString)
	}

	return t.DB.GetTasks()
}

func (t *TaskUsecase) GetTaskById(id string) (*model.Task, error) {
	return t.DB.GetTaskById(id)
}

func (t *TaskUsecase) UpdateTask(task *model.Task, pastDay bool) error {
	_, err := t.DB.GetTaskById(task.Id)
	if err != nil {
		return err
	}

	if pastDay {
		nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = nextDate
	}

	return t.DB.UpdateTask(task)
}

func (t *TaskUsecase) MakeTaskDone(id string) error {
	task, err := t.DB.GetTaskById(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return t.DB.DeleteTask(id)
	}

	nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return t.DB.MakeTaskDone(id, nextDate)
}

func (t *TaskUsecase) DeleteTask(id string) error {
	return t.DB.DeleteTask(id)
}

func parseValue(num string) (int, error) {
	days, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}

	if days >= 400 || days < 0 {
		return 0, fmt.Errorf("invalid value %d", days)
	}

	return days, nil
}

func addDateTask(now time.Time, dateTask time.Time, year int, month int, day int) time.Time {
	dateTask = dateTask.AddDate(year, month, day)

	for dateTask.Before(now) {
		dateTask = dateTask.AddDate(year, month, day)
	}

	return dateTask
}

func getDateTaskByWeek(now, dateTask time.Time, daysString string) (time.Time, error) {
	daysSlice := strings.Split(daysString, ",")
	daysOfWeek := regexp.MustCompile("[1-7]")

	daysOfWeekMap := make(map[int]bool)
	for _, day := range daysSlice {
		numberOfDay, err := strconv.Atoi(day)
		if err != nil {
			return dateTask, err
		}

		if len(daysOfWeek.FindAllString(day, -1)) == 0 {
			return dateTask, fmt.Errorf("invalid value %d day of the week", numberOfDay)
		}

		if numberOfDay == sundayRus {
			numberOfDay = sundayEU
		}
		daysOfWeekMap[numberOfDay] = true
	}

	for {
		if daysOfWeekMap[int(dateTask.Weekday())] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1)
	}

	return dateTask, nil
}

func getDateTaskByMonth(now, dateTask time.Time, repeat []string) (time.Time, error) {
	daysString := repeat[1]

	monthsString := ""
	if len(repeat) > 2 {
		monthsString = repeat[2]
	}

	daysSlice := strings.Split(daysString, ",")
	monthsSlice := strings.Split(monthsString, ",")

	daysMap := make(map[int]bool)
	for _, day := range daysSlice {
		numberOfDay, err := strconv.Atoi(day)
		if err != nil {
			return dateTask, err
		}
		if numberOfDay < -2 || numberOfDay > 31 || numberOfDay == 0 {
			return dateTask, fmt.Errorf("invalid value %d day of the month", numberOfDay)
		}
		daysMap[numberOfDay] = true
	}

	monthsMap := make(map[int]bool)
	for _, month := range monthsSlice {
		if month == "" {
			continue
		}
		numberOfMonth, err := strconv.Atoi(month)
		if err != nil {
			return dateTask, err
		}
		if numberOfMonth < 1 || numberOfMonth > 12 {
			return dateTask, fmt.Errorf("invalid value %d month", numberOfMonth)
		}
		monthsMap[numberOfMonth] = true
	}

	for {
		if len(monthsMap) == 0 {
			break
		}

		if monthsMap[int(dateTask.Month())] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1)
	}

	for {
		lastDayOfMonth := time.Date(dateTask.Year(), dateTask.Month()+1, 0, 0, 0, 0, 0, dateTask.Location()).Day()
		predLastDayOfMonth := lastDayOfMonth - 1
		fmt.Println(predLastDayOfMonth)

		key := dateTask.Day()
		switch {
		case lastDayOfMonth == dateTask.Day():
			if _, ok := daysMap[-1]; ok {
				key = -1
			}
		case predLastDayOfMonth == dateTask.Day():
			if _, ok := daysMap[-2]; ok {
				key = -2
			}
		}

		if daysMap[key] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1)
	}

	return dateTask, nil
}
