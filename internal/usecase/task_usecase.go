package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"todo/internal/entity"
)

type repository interface {
	// AddTask adds a new task to the SQLite database.
	// It takes a context, representing the execution context, and a task, representing the task to be added.
	// It returns the ID of the newly added task or an error if any.
	AddTask(ctx context.Context, task *entity.Task) (int64, error)

	// GetTasks retrieves tasks from the SQLite database.
	// It takes a context, representing the execution context,
	// and a limit, representing the maximum number of tasks to retrieve.
	// It returns the retrieved tasks or an error if any.
	GetTasks(ctx context.Context, limit int) ([]*entity.Task, error)

	// GetTasksByDate retrieves tasks from the SQLite database by date.
	// It takes a context, representing the execution context, and a date, representing the date to filter tasks by.
	// It returns the retrieved tasks or an error if any.
	GetTasksByDate(ctx context.Context, date string) ([]*entity.Task, error)

	// GetTasksByQuery retrieves tasks from the SQLite database by a search query.
	// It takes a context, representing the execution context, and a query, representing the search query.
	// It returns the retrieved tasks or an error if any.
	GetTasksByQuery(ctx context.Context, query string) ([]*entity.Task, error)

	// DeleteTask deletes a task from the SQLite database by its ID.
	// It takes a context, representing the execution context, and an ID, representing the ID of the task to delete.
	// It returns an error if any.
	DeleteTask(ctx context.Context, id int64) error

	// GetTask retrieves a task from the SQLite database by its ID.
	// It takes a context, representing the execution context, and an ID, representing the ID of the task to retrieve.
	// It returns the retrieved task or an error if any.
	GetTask(ctx context.Context, id int64) (*entity.Task, error)

	// UpdateTask updates a task in the SQLite database.
	// It takes a context, representing the execution context, and a task, representing the updated task.
	// It returns an error if any.
	UpdateTask(ctx context.Context, task *entity.Task) error
}

type jwt interface {
	// CreateToken creates a new JWT token using the provided secret key.
	// It returns the token string and any error encountered during the process.
	CreateToken() (string, error)

	// ValidateToken validates the given JWT token string.
	// It parses the token using the secret key stored in the JWT instance.
	// If the token is valid, it returns nil. Otherwise, it returns an error.
	ValidateToken(tokenString string) error
}

// TaskUsecase represents the use case for managing tasks.
type TaskUsecase struct {
	repo     repository
	jwt      jwt
	password string
}

// NewTaskUsecase creates a new instance of the TaskUsecase.
// It takes a repository as a parameter and an optional list of options.
// The options are applied to the TaskUsecase using the functional option pattern.
// Returns a pointer to the created TaskUsecase.
func NewTaskUsecase(repo repository, jwt jwt, opts ...Option) *TaskUsecase {
	tu := &TaskUsecase{
		repo: repo,
		jwt:  jwt,
	}

	for _, opt := range opts {
		opt(tu)
	}

	return tu
}

// Login authenticates the user with the provided password.
// It returns a string representing the user's authentication token and an error if any.
func (tu *TaskUsecase) Login(password string) (string, error) {
	if tu.password == "" {
		return "", entity.ErrAuthDisabled
	}

	if tu.password != password {
		return "", entity.ErrInvalidPassword
	}

	token, err := tu.jwt.CreateToken()
	if err != nil {
		return "", fmt.Errorf("jwt.CreateToken: %w", err)
	}

	return token, nil
}

// ValidateToken validates the given token using JWT validation.
// It returns an error if the token is invalid or if JWT validation fails.
func (tu *TaskUsecase) ValidateToken(token string) error {
	if tu.password == "" {
		return entity.ErrAuthDisabled
	}

	if err := tu.jwt.ValidateToken(token); err != nil {
		return fmt.Errorf("jwt.ValidateToken: %w", err)
	}

	return nil
}

// ShouldCheckToken checks if the password is set
// and returns a boolean value indicating whether token should be checked.
func (tu *TaskUsecase) ShouldCheckToken() bool {
	return tu.password != ""
}

// AddTask adds a new task to the repository.
// It takes a context, representing the execution context, and a task, representing the task to be added.
// It returns the ID of the newly added task or an error if any.
func (tu *TaskUsecase) AddTask(ctx context.Context, task *entity.Task) (string, error) {
	var err error

	if task.Title == "" {
		return "", entity.ErrEmptyTitle
	}

	if task.Date == "" {
		task.Date = time.Now().Truncate(entity.DayInHours * time.Hour).Format(entity.DateFormat)
	}

	date, err := time.Parse(entity.DateFormat, task.Date)
	if err != nil {
		return "", fmt.Errorf("time.Parse: %w", err)
	}

	var nextDate string
	if task.Repeat == "" {
		nextDate = time.Now().Format(entity.DateFormat)
	} else {
		nextDate, err = calculateNextDate(time.Now().Truncate(entity.DayInHours*time.Hour), task.Date, task.Repeat)
		if err != nil {
			return "", fmt.Errorf("calculateNextDate: %w", err)
		}
	}

	if date.Before(time.Now().Truncate(entity.DayInHours * time.Hour)) {
		task.Date = nextDate
	}

	id, err := tu.repo.AddTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("repo.AddTask: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}

// GetTasks retrieves a list of tasks from the repository.
// It returns a slice of Task entities and an error, if any.
func (tu *TaskUsecase) GetTasks(ctx context.Context, query string) ([]*entity.Task, error) {
	var (
		err   error
		tasks []*entity.Task
	)

	if query == "" {
		tasks, err = tu.repo.GetTasks(ctx, entity.LimitTasks)
		if err != nil {
			return nil, fmt.Errorf("repo.GetTasks: %w", err)
		}

		return tasks, nil
	}

	date, err := time.Parse("02.01.2006", query)
	if err == nil {
		tasks, err = tu.repo.GetTasksByDate(ctx, date.Format(entity.DateFormat))
		if err != nil {
			return nil, fmt.Errorf("repo.GetTasksByDate: %w", err)
		}

		return tasks, nil
	}

	tasks, err = tu.repo.GetTasksByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo.GetTasksByQuery: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a task by its ID.
// It takes a context and an ID as input parameters and returns the corresponding task and any error encountered.
func (tu *TaskUsecase) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	task, err := tu.repo.GetTask(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetTask: %w", err)
	}

	return task, nil
}

// GetNextDate calculates the next date based on the current date, a given date, and a repeat pattern.
// It takes the current date in the format specified by the entity.DateFormat constant,
// and returns the next date in the same format.
// If any error occurs during the calculation, it returns an empty string and the error.
func (tu *TaskUsecase) GetNextDate(now, date, repeat string) (string, error) {
	timeNow, err := time.Parse(entity.DateFormat, now)
	if err != nil {
		return "", fmt.Errorf("time.Parse: %w", err)
	}

	nextDate, err := calculateNextDate(timeNow, date, repeat)
	if err != nil {
		return "", fmt.Errorf("calculateNextDate: %w", err)
	}

	return nextDate, nil
}

// DeleteTask deletes a task with the specified ID.
// It first checks if the task exists by calling the GetTask method of the repository.
// If the task exists, it then calls the DeleteTask method of the repository to delete the task.
// If any error occurs during the process, it returns an error with additional context information.
func (tu *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}

	if _, err = tu.repo.GetTask(ctx, parsedID); err != nil {
		return fmt.Errorf("repo.GetTask: %w", err)
	}

	if err = tu.repo.DeleteTask(ctx, parsedID); err != nil {
		return fmt.Errorf("repo.DeleteTask: %w", err)
	}

	return nil
}

// UpdateTask updates the given task in the repository.
// It validates the task's ID, title, and date, and sets the date to the current time if it is empty.
// It also calculates the next date based on the repeat interval, if provided.
// If the calculated date is in the past, it updates the task's date to the next date.
// Finally, it calls the repository's UpdateTask method to update the task.
// If any validation or repository operation fails, it returns an error.
func (tu *TaskUsecase) UpdateTask(ctx context.Context, task *entity.Task) error {
	if task.ID == "" {
		return entity.ErrEmptyID
	}

	parsedID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}

	if _, err = tu.repo.GetTask(ctx, parsedID); err != nil {
		return entity.ErrTaskNotFound
	}

	if task.Title == "" {
		return entity.ErrEmptyTitle
	}

	if _, err = time.Parse(entity.DateFormat, task.Date); err != nil {
		return entity.ErrInvalidDate
	}

	if _, err = calculateNextDate(time.Now(), task.Date, task.Repeat); task.Repeat != "" && err != nil {
		return err
	}

	if err = tu.repo.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("repo.UpdateTask: %w", err)
	}

	return nil
}

// DoTask performs a task based on the given task ID.
// If the task has a repeat schedule, it calculates the next date and updates the task.
// If the task does not have a repeat schedule, it deletes the task from the repository.
// Returns an error if any operation fails.
func (tu *TaskUsecase) DoTask(ctx context.Context, id string) error {
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}

	task, err := tu.repo.GetTask(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("repo.GetTask: %w", err)
	}

	if task.Repeat == "" {
		if err = tu.repo.DeleteTask(ctx, parsedID); err != nil {
			return fmt.Errorf("repo.DeleteTask: %w", err)
		}

		return nil
	}

	nextDate, err := calculateNextDate(time.Now().Truncate(entity.DayInHours*time.Hour), task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("calculateNextDate: %w", err)
	}

	task.Date = nextDate

	if err = tu.repo.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("repo.UpdateTask: %w", err)
	}

	return nil
}

func calculateNextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", entity.ErrMissingRepeatParams
	}

	startDate, err := time.Parse(entity.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("time.Parse: %w", err)
	}

	parts := strings.Split(repeat, " ")
	param := parts[0]
	switch param {
	case "y":
		return calculateNextDateYear(now, startDate)
	case "d":
		return calculateNextDateDay(now, startDate, parts)
	case "w":
		return calculateNextDateWeek(now, parts)
	case "m":
		return calculateNextDateMonth(now, startDate, parts)
	default:
		return "", entity.ErrUnsupportedRepeatFormat
	}
}

func calculateNextDateYear(now, startDate time.Time) (string, error) {
	currDate := startDate.AddDate(1, 0, 0)
	for now.After(currDate) || now.Equal(currDate) {
		currDate = currDate.AddDate(1, 0, 0)
	}

	return currDate.Format(entity.DateFormat), nil
}

func calculateNextDateDay(now, startDate time.Time, parts []string) (string, error) {
	if len(parts) == 1 {
		return "", entity.ErrNoInterval
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("strconv.Atoi: %w", err)
	}

	if days > entity.MaxIntervalInDays {
		return "", entity.ErrMaxIntervalExceeded
	}

	currDate := startDate.AddDate(0, 0, days)
	for now.After(currDate) {
		currDate = currDate.AddDate(0, 0, days)
	}

	return currDate.Format(entity.DateFormat), nil
}

func calculateNextDateWeek(now time.Time, parts []string) (string, error) {
	var weekdays []int
	if len(parts) == 1 {
		return "", entity.ErrNoInterval
	}

	weekdays, err := extIntParams(parts[1], 1, entity.DaysInWeek)
	if err != nil {
		return "", err
	}

	currDate := now.AddDate(0, 0, 1)
	for {
		day := int(currDate.Weekday())
		for _, weekday := range weekdays {
			if day == weekday || (day == 0 && weekday == 7) {
				return currDate.Format(entity.DateFormat), nil
			}
		}
		currDate = currDate.AddDate(0, 0, 1)
	}
}

func calculateNextDateMonth(now, startDate time.Time, parts []string) (string, error) {
	if len(parts) == 1 {
		return "", entity.ErrNoInterval
	}

	monthDays, err := extIntParams(parts[1], -2, entity.MaxDaysInMonth)
	if err != nil {
		return "", err
	}

	sortDayParams(monthDays)

	const customIntervalParts = 3
	var months []int
	if len(parts) == customIntervalParts {
		months, err = extIntParams(parts[2], 1, entity.MonthsInYear)
		if err != nil {
			return "", err
		}

		sort.Ints(months)
	} else {
		s := int(startDate.Month())
		n := int(now.Month())
		if startDate.After(now) {
			months = []int{s, s + 1}
		} else {
			months = []int{n, n + 1}
		}
	}

	dateMap := buildDateMap(months)
	for _, d := range monthDays {
		for _, m := range months {
			days := dateMap[m]
			if len(days) <= d-1 {
				continue
			}

			var currDate time.Time
			if d < 0 {
				currDate = days[len(days)+d]
			} else {
				currDate = days[d-1]
			}

			if now.Before(currDate) || now.Equal(currDate) {
				return currDate.Format(entity.DateFormat), nil
			}
		}
	}

	return "", entity.ErrDateCalculation
}

func buildDateMap(months []int) map[int][]time.Time {
	res := map[int][]time.Time{}
	for _, m := range months {
		currDay := time.Date(time.Now().Year(), time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		lastDay := currDay.AddDate(0, 1, -1)
		res[m] = make([]time.Time, lastDay.Day())
		for i := range lastDay.Day() {
			res[m][i] = currDay
			currDay = currDay.AddDate(0, 0, 1)
		}
	}
	return res
}

func extIntParams(params string, min, max int) ([]int, error) {
	strings := strings.Split(params, ",")

	numbers := make([]int, len(strings))
	for i := range len(strings) {
		number, err := strconv.Atoi(strings[i])
		if err != nil {
			return nil, err
		}
		if number < min || number > max {
			return nil, fmt.Errorf("illegal value %d", number)
		}
		numbers[i] = number
	}

	return numbers, nil
}

func sortDayParams(days []int) {
	sort.SliceStable(days, func(i, j int) bool {
		if days[i] < 0 && days[j] < 0 {
			return days[i] < days[j]
		}
		if days[i] < 0 {
			return false
		}
		if days[j] < 0 {
			return true
		}
		return days[i] < days[j]
	})
}
