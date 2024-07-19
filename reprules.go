package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Функция для вычисления следующей даты на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	fmt.Println("======= Принятые значения: Сейчас:", now.Format("20060102"), "--Старт:", date, "--", repeat, "===========")

	// Проверяем наличие правила повторения, если его нет - возвращаем ошибку
	if repeat == "" {
		fmt.Println("======= Правило повторения отсутствует! ===========")
		return "", errors.New("правило повторения отсутствует")
	}

	rep := strings.Split(repeat, " ")
	fmt.Println("======= Парсим правило повторения:", rep, "===========")

	if len(rep) < 1 || (rep[0] != "y" && rep[0] != "d") {
		fmt.Println("======= Неподдерживаемое правило повторения! ===========")
		return "", errors.New("неподдерживаемый формат правила повторения")
	}

	// Парсим дату события в формате YYYYMMDD
	timBase, err := time.Parse("20060102", date)
	if err != nil {
		fmt.Println("======= Некорректная дата! ===========")
		return "", err
	}
	fmt.Println("======= Дата успешно распознана:", timBase, "===========")

	// Проверяем режим повторения: год или день
	if rep[0] == "y" {
		// Если год, то прибавляем к дате год, пока не найдем следующую дату после текущей
		fmt.Println("======= Определяем режим повтора: год (y) ===========")
		for timBase.Before(now) {
			timBase = timBase.AddDate(1, 0, 0) // Добавляем один год
			fmt.Println("======= Добавляем 1 год! Новая дата:", timBase.Format("20060102"), "===========")
		}
		result := timBase.Format("20060102")
		fmt.Println("======= Старая дата:", date, "===========")
		fmt.Println("=======  Новая дата:", result, "===========")
		return result, nil
	}

	if rep[0] == "d" {
		// Если день, то прибавляем указанное количество дней
		fmt.Println("======= Определяем режим повтора: день (d) ===========")
		if len(rep) < 2 {
			fmt.Println("======= Некорректно указан режим повторения! ===========")
			return "", errors.New("некорректно указан режим повторения")
		}

		days, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", err // Возвращаем ошибку, если количество дней некорректно
		}

		if days > 400 {
			fmt.Println("======= Количество дней превышает 400! ===========")
			return "", errors.New("перенос события более чем на 400 дней недопустим")
		}

		fmt.Println("======= Количество дней для добавления:", days, "===========")
		for timBase.Before(now) {
			timBase = timBase.AddDate(0, 0, days)
			fmt.Println("======= Добавляем", days, "дней! Новая дата:", timBase.Format("20060102"), "===========")
		}

		result := timBase.Format("20060102")
		fmt.Println("=======   Текущая дата:", now.Format("20060102"), "===========")
		fmt.Println("======= Стартовая дата:", date, "===========")
		fmt.Println("=======     Новая дата:", result, "===========")
		return result, nil
	}

	return "", errors.New("некорректное правило повторения")
}
