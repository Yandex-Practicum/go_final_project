package utils

import "errors"

const (
	DefaultPort = 7540
	EnvPassword = "TODO_PASSWORD"
	EnvPort     = "TODO_PORT"
	WebDir      = "./web"
)

const (
	AuthSecret = "kyusugfi7b6234udb7o38268bd5bgk23uk"
)

const (
	ErrUnsupportedMethod    = "Метод не поддерживается"
	ErrInvalidDateNowFormat = "Invalid 'now' date format"
)

const (
	ParseDateFormat = "20060102"
)

var (
	ErrIDIsEmpty         = errors.New("не указан идентификатор")
	ErrInvalidTaskTitle  = errors.New("не указан заголовок задачи")
	ErrInvalidJson       = errors.New("ошибка десериализации JSON")
	ErrInvalidTaskDate   = errors.New("неверно указана дата задачи")
	ErrInvalidTaskRepeat = errors.New("неверно указана дата задачи и повтор")
	ErrTaskParse         = errors.New("ошибка разбора задач из базы данных")
	ErrTaskNotFound      = errors.New("задача не найдена")
	ErrDBInsert          = errors.New("ошибка вставки в базу данных")
	ErrGetTaskID         = errors.New("ошибка получения ID задачи")
	ErrInvalidPassword   = errors.New("не верный пароль")
	ErrTokenCreate       = errors.New("ошибка генерации токена")
)

const (
	FilterTypeNone   = iota
	FilterTypeDate   = iota
	FilterTypeSearch = iota
)
