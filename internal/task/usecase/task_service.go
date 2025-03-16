package usecase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"test/internal/task"
	"time"
)

type RepositoryInterface interface {
	Insert(task *task.Task) (string, error)
	GetAll() (*task.List, error)
	GetByDate(date string) (*task.List, error)
	GetByTitleOrComment(search string) (*task.List, error)
	GetById(id int) (*task.Task, error)
	DeleteById(id int) error
	UpdateById(t *task.Task) (*task.Task, error)
	Done(task *task.Task) error
}

type TaskService struct {
	taskRepository RepositoryInterface
}

func NewTaskService(taskRepository RepositoryInterface) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (ts *TaskService) NextDate(now string, date string, repeat string) (error, string) {
	if repeat == "" {
		return errors.New("правило пустое"), ""
	}

	startDate, err := time.Parse(task.FormatDate, date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %v", err), ""
	}

	nowDate, err := time.Parse(task.FormatDate, now)
	if err != nil {
		return fmt.Errorf("неверный формат даты %v", err), ""
	}

	parts := strings.Fields(repeat)
	if len(parts) < 1 {
		return errors.New("неверный формат правила повторения"), ""
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return errors.New("неверный формат правила d"), ""
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 365 {
			return errors.New("неверный интервал"), ""
		}
		for {
			startDate = startDate.AddDate(0, 0, interval)
			if startDate.After(nowDate) {
				break
			}
		}
	case "w":
		if len(parts) != 2 {
			return errors.New("неверный формат правила w"), ""
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 7 {
			return errors.New("неверный интервал недели"), ""
		}
		for {
			startDate = startDate.AddDate(0, 0, 1)
			if int(startDate.Weekday()) == interval && startDate.After(nowDate) {
				break
			}
		}
	case "m":
		if len(parts) != 2 {
			return errors.New("неверный формат m"), ""
		}
		dayParts := strings.Split(parts[1], ",")
		if len(dayParts) != 2 {
			return errors.New("неверные дни месяцев"), ""
		}
		monthOffset, err := strconv.Atoi(dayParts[0])
		if err != nil {
			return errors.New("ошибка конвертации"), ""
		}
		dayOfMonth, err := strconv.Atoi(dayParts[1])
		if err != nil || dayOfMonth < 1 || dayOfMonth > 31 {
			return errors.New("неверный день месяца"), ""
		}
		for {
			startDate = startDate.AddDate(0, monthOffset, 0)
			if startDate.Day() != dayOfMonth {
				startDate = time.Date(startDate.Year(), startDate.Month(), dayOfMonth, 0, 0, 0, 0, startDate.Location())
			}
			if startDate.After(nowDate) {
				break
			}
		}
	case "y":
		if len(parts) != 1 {
			return errors.New("неверный формат правила y"), ""
		}
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(nowDate) {
				break
			}
		}
	default:
		return errors.New("неизвестный тип правила"), ""
	}

	return nil, startDate.Format(task.FormatDate)
}

func (ts *TaskService) Create(t *task.Task) (string, error) {
	dateNow := time.Now().Format(task.FormatDate)
	err := t.ValidateForCreate(dateNow)
	if err != nil {
		return "", err
	}

	if t.Repeat == "" || t.Date == dateNow {
		t.Date = dateNow
	} else {
		err, t.Date = ts.NextDate(dateNow, t.Date, t.Repeat)
		if err != nil {
			return "", err
		}
	}

	id, err := ts.taskRepository.Insert(t)
	if err != nil {
		return "", err
	}

	t.ID = id

	return id, nil
}

func (ts *TaskService) GetAll(search string) (*task.List, error) {
	if search != "" {
		_, err := time.Parse("02.01.2006", search)
		if err == nil {
			return ts.taskRepository.GetByDate(search)
		}

		return ts.taskRepository.GetByTitleOrComment(search)
	}

	taskList, err := ts.taskRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (ts *TaskService) GetById(id int) (*task.Task, error) {
	t, err := ts.taskRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskService) Update(t *task.Task) (*task.Task, error) {
	dateNow := time.Now().Format(task.FormatDate)
	err := t.ValidateForCreate(dateNow)
	if err != nil {
		return nil, err
	}

	err, t.Date = ts.NextDate(dateNow, t.Date, t.Repeat)
	if err != nil {
		return nil, err
	}

	t, err = ts.taskRepository.UpdateById(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (ts *TaskService) Delete(id int) error {
	err := ts.taskRepository.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) Done(paramId string) error {
	if paramId == "" {
		return errors.New("задача не найдена")
	}
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return errors.New("некорректный параметр id")
	}
	t, err := ts.GetById(id)
	if err != nil {
		return errors.New("задача не найдена")
	}
	if t.Repeat == "" {
		err = ts.Delete(id)
		if err != nil {
			return errors.New("ошибка при удалении")
		}

		return nil
	}

	err, newDate := ts.NextDate(time.Now().Format(task.FormatDate), t.Date, t.Repeat)
	if err != nil {
		return errors.New("ошибка при вычислении даты")
	}
	t.Date = newDate
	err = ts.taskRepository.Done(t)
	if err != nil {
		return err
	}

	return nil
}
