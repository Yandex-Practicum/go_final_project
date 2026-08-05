package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const LayoutDate string = "20060102"

func mounthCheck(day, mounth string) (string, bool) { // day or mounth
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		log.Println("Ошибка перевода строки в целое число в функции проверки дня месяца")
		return "", false
	}
	if dayInt < -2 || dayInt > 31 || dayInt == 0 {
		return "", false
	}

	mounthInt, err := strconv.Atoi(mounth)
	if err != nil {
		log.Println("Ошибка перевода строки в целое число в функции проверки дня месяца")
		return "", false
	}
	if mounthInt < 1 || mounthInt > 12 {
		return "", false
	}

	res := fmt.Sprintln(day, mounth)
	return res, true
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(LayoutDate, dstart)
	if err != nil {
		log.Printf("Ошибка парсинга даты в тип Тайма: %v", err)
		return "", err
	}

	if repeat == "" || len(repeat) < 1 {
		return "", fmt.Errorf("repeat не может быть пустым")
	}

	rep := strings.Split(repeat, " ")
	if rep[0] != "d" && rep[0] != "y" && rep[0] != "w" && rep[0] != "m" {
		return "", fmt.Errorf("Ошибка при проверке repeat, первый символ не является допустимым.")
	}

	switch rep[0] {

	case "m":
		if len(rep) != 2 && len(rep) != 3 {
			return "", fmt.Errorf("repeat  не может содержать меньше 1 и больше 3 параметров")
		}
		var daysAll [34]bool
		days := strings.Split(rep[1], ",")
		for _, dayMou := range days {
			day, err := strconv.Atoi(dayMou)
			if err != nil {
				return "", fmt.Errorf("день месяца должен быть числом: %w", err)
			}
			switch {
			case day >= 1 && day <= 31:
				daysAll[day] = true
			case day == -1:
				daysAll[32] = true
			case day == -2:
				daysAll[33] = true
			default:
				return "", fmt.Errorf("Некорректный день месяца: %w", err)
			}
		}
		var mounthAll [13]bool
		if len(rep) == 2 {
			for mounth := 1; mounth <= 12; mounth++ {
				mounthAll[mounth] = true
			}
		} else {
			for _, v := range strings.Split(rep[2], ",") {
				mounth, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("месяц должен быть числом: %w", err)
				}
				if mounth < 1 || mounth > 12 {
					return "", fmt.Errorf(" месяц должен быть от 1 до 12")
				}
				mounthAll[mounth] = true
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if !date.After(now) {
				continue
			}
			day := date.Day()
			month := int(date.Month())
			if !mounthAll[month] {
				continue
			}
			isLastDay := date.AddDate(0, 0, 1).Day() == 1
			isPenultimateDay := date.AddDate(0, 0, 2).Day() == 1
			if daysAll[day] || (daysAll[32] && isLastDay) || (daysAll[33] && isPenultimateDay) {
				break
			}
		}
	case "w":
		if len(rep) != 2 {
			return "", fmt.Errorf("repeat w не может содержать не 2 параметра")
		}
		var allowed [8]bool
		days := strings.Split(rep[1], ",")
		for _, dayStr := range days {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", fmt.Errorf("день недели должен быть числом %w", err)
			}
			if day < 1 || day > 7 {
				return "", fmt.Errorf("Передан невалидный день недели, нужно от 1 до 7")
			}
			allowed[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if date.After(now) && allowed[weekday] {
				break
			}
		}

	case "d":

		if len(rep) != 2 {
			return "", fmt.Errorf("Для повтора d нужен шаг: например, d 7")
		}

		step, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", fmt.Errorf("Некорректный шаг дней: %w", err)
		}
		if step < 1 || step > 400 {
			return "", fmt.Errorf("Шаг дней должен быть от 1 до 400")
		}

		for {
			date = date.AddDate(0, 0, step)
			if date.After(now) {
				break
			}
		}

	case "y":
		if len(rep) != 1 {
			return "", fmt.Errorf("Для ежегодного повтора необходимо предоставить [y]")
		}

		for {
			date = date.AddDate(1, 0, 0)

			if date.After(now) {
				break
			}
		}
	}

	res := date.Format(LayoutDate)
	return res, nil

}
