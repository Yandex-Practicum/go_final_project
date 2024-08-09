package utils

import "errors"

const (
	EnvPassword = "TODO_PASSWORD"
)

const (
	AuthSecret = "kyusugfi7b6234udb7o38268bd5bgk23uk"
)

const (
	ErrUnsupportedMethod    = "Метод не поддерживается"
	ErrInvalidJson          = "Ошибка десериализации JSON"
	ErrInvalidDateNowFormat = "Invalid 'now' date format"
	ErrIDIsEmpty            = "Не указан идентификатор"
	ErrTaskNotFound         = "Задача не найдена"
	ErrTaskParse            = "Ошибка разбора задач из базы данных"
	ErrDBInsert             = "Ошибка вставки в базу данных"
	ErrGetTaskID            = "Ошибка получения ID задачи"
)

const (
	ParseDateFormat = "20060102"
)

var (
	ErrInvalidTaskTitle  = errors.New("не указан заголовок задачи")
	ErrInvalidTaskDate   = errors.New("неверно указана дата задачи")
	ErrInvalidTaskRepeat = errors.New("неверно указана дата задачи и повтор")
)
