package constants

import "errors"

const (
	DefaultPort = 7540
	EnvPassword = "TODO_PASSWORD"
	EnvPort     = "TODO_PORT"
	EnvSecret   = "AUTH_SECRET"
	WebDir      = "./web"
)

const (
	ErrUnsupportedMethod    = "Метод не поддерживается"
	ErrInvalidDateNowFormat = "Invalid 'now' date format"
)

const (
	ParseDateFormat = "20060102"
)

var (
	ErrIdIsEmpty         = errors.New("не указан идентификатор")
	ErrInvalidTaskTitle  = errors.New("не указан заголовок задачи")
	ErrInvalidJson       = errors.New("ошибка десериализации JSON")
	ErrInvalidTaskDate   = errors.New("неверно указана дата задачи")
	ErrInvalidTaskRepeat = errors.New("неверно указана дата задачи и повтор")
	ErrTaskParse         = errors.New("ошибка разбора задач из базы данных")
	ErrTaskNotFound      = errors.New("задача не найдена")
	ErrDBInsert          = errors.New("ошибка вставки в базу данных")
	ErrGetTaskId         = errors.New("ошибка получения Id задачи")
	ErrInvalidPassword   = errors.New("неверный пароль")
	ErrTokenCreate       = errors.New("ошибка генерации токена")
)

const (
	FilterTypeNone   = iota
	FilterTypeDate   = iota
	FilterTypeSearch = iota
)
