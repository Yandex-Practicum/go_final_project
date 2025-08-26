package handler

import "github.com/askoren1/go_final_project/internal/repository"

//Определяем структуру Handler и функцию-конструктор New для создания экземпляров этой структуры.
//Структура Handler используется для объединения логики работы с репозиторием (базой данных) и HTTP-обработчиками

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler { //функция-конструктор для типа Handler
	return &Handler{
		repo: repo,
	}
}
