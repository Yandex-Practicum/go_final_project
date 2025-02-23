package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const f = "20060102"

// ф-я для M. если третий параметр сущ-т && в нем есть текущий месяц || третьего параметра нет - дата подходит
func shouldReturnDate(month int, monthSlice []int) bool {
	if monthSlice == nil {
		return true
	}
	for _, m := range monthSlice {
		if m == month {
			return true
		}
	}
	return false
}

// ф-я для M. переводит слайс строк в слайс чисел
func SliceStrToIntM(SecondStr []string, min int, max int) ([]int, error) {

	intSlice := make([]int, len(SecondStr))

	for i, day := range SecondStr {
		d, err := strconv.Atoi(day)
		if err != nil || (d < min || d > max) {
			return nil, fmt.Errorf("нужно ук-ть числа в пределах от %d, до%d", min, max)
		}
		intSlice[i] = d
	}
	return intSlice, nil
}

// ф-я для W. переводит слайс строк в слайс чисел
func SliceStrToIntW(SecondStr []string) ([]int, error) {

	intSlice := make([]int, len(SecondStr))

	for i, day := range SecondStr {
		d, err := strconv.Atoi(day)
		if err != nil || (d < 1 || d > 7) {
			return nil, fmt.Errorf("нужно указать дни недели в числовом формате (1-7): %w", err)
		}
		intSlice[i] = d
		if i > 0 {
			if intSlice[i] <= intSlice[i-1] {
				return nil, fmt.Errorf("необходимые дни недели ук-ся в порядке возрастания: %w", err)
			}
		}
	}
	return intSlice, nil
}

// ф-я для M. Находит последний день месяца
func GetLastDayOfMonth(year int, month time.Month) time.Time {
	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfNextMonth.AddDate(0, 0, -1)
	return lastDayOfMonth
}

// основная ф-я. Обрабатывает все правила
func NextDate(nowTime string, date string, repeat string) (string, error) {

	now, err := time.Parse(f, nowTime) //парсинг даты
	if err != nil {
		return "", fmt.Errorf("неверный формат даты 'now'")
	}

	d, err := time.Parse(f, date) //парсинг даты запроса
	if err != nil {
		return "", fmt.Errorf("неверный формат даты 'date'")
	}
	if repeat == "" {
		return "", fmt.Errorf("") //вроде как добавить ф-ю, удаляющую задачу
	}

	parts := strings.Split(repeat, " ")

	letter := ""
	if len(parts) != 0 {
		letter = parts[0]
	}

	if letter != "d" && letter != "y" && letter != "w" && letter != "m" {
		return "", fmt.Errorf("неправильная буква")
	}
	if len(parts) > 3 {
		return "", fmt.Errorf("лишние данные 'date'")
	}
	// правила для y
	if letter == "y" {
		if len(parts) == 1 {
			var NewDate time.Time
			NewDate = d
			for {
				NewDate = NewDate.AddDate(1, 0, 0)
				if NewDate.Format(f) > now.Format(f) {
					break
				}
			}
			return NewDate.Format(f), nil
		} else {
			return "", fmt.Errorf("лишние параметры")
		}
	}

	SecondPart := ""
	if len(parts) >= 2 {
		SecondPart = parts[1]
	}
	SecondStr := strings.Split(SecondPart, ",")

	// правила для d
	if letter == "d" {
		if len(parts) == 1 {
			return "", fmt.Errorf("не указан интервал в днях")
		}
		if len(parts) == 2 {
			if len(SecondStr) == 1 {
				SecondInt, err := strconv.Atoi(SecondStr[0])
				if err != nil {
					return "", fmt.Errorf("второй параметр должен быть числом")
				}
				if SecondInt >= 1 && SecondInt <= 400 {
					var NewDate time.Time
					NewDate = d
					for {
						NewDate = NewDate.AddDate(0, 0, SecondInt)

						if NewDate.Format(f) > now.Format(f) {
							break
						}

					}
					return NewDate.Format(f), nil
				} else {
					return "", fmt.Errorf("недопустимое количество дней")
				}
			}
		} else {
			return "", fmt.Errorf("лишние параметры")
		}
	}

	// правила для w
	if letter == "w" {
		if len(parts) == 1 {
			return "", fmt.Errorf("не указан день недели")
		}
		if len(parts) == 2 {

			//перевожу текстовый слайс в числовой формат
			intSlice, err := SliceStrToIntW(SecondStr)
			if err != nil {
				return "", err
			}

			// Получаю число дня недели и преобразую в соотв-и пн=1,вс=7
			nowDay := int(d.Weekday())
			if nowDay == 0 {
				nowDay = 7 // Воскресенье (0) становится 7
			}

			smallestDay := intSlice[0]
			thisWeekDay := intSlice[0]

			//смотрю, есть ли в слайсе день недели больше сегоднящнего
			for _, weekday := range intSlice {
				if weekday > nowDay {
					thisWeekDay = weekday
					break
				}
			}
			//если есть, получаю новую дату с thisWeekDay и сравниваю с now
			if thisWeekDay > nowDay {

				thisWeekDay = thisWeekDay - nowDay
				newDate := time.Date(d.Year(), d.Month(), d.Day()+thisWeekDay, 0, 0, 0, 0, d.Location())
				log.Printf("newDate создан")
				if newDate.Format(f) > now.Format(f) {

					return newDate.Format(f), nil
				} else {
				}

			}

			smallestDay = (7 - nowDay) + smallestDay
			newDate := time.Date(d.Year(), d.Month(), d.Day()+smallestDay, 0, 0, 0, 0, d.Location())
			for {
				if newDate.Format(f) > now.Format(f) {
					return newDate.Format(f), nil
				}
				newDate = time.Date(newDate.Year(), newDate.Month(), newDate.Day()+7, 0, 0, 0, 0, newDate.Location())

			}
		}
		return "", fmt.Errorf("лишние параметры")
	}

	// правила для m (m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12])
	if letter == "m" {
		if len(parts) == 1 {
			return "", fmt.Errorf("необходимо указать день месяца")
		}
		if len(parts) <= 3 {

			//заранее создаю слайс для третьего параметра, чтобы обозначить его видимость
			var monthSlice []int

			if len(parts) == 3 {
				ThirdPart := parts[2]
				ThirdStr := strings.Split(ThirdPart, ",")

				monthSlice, err = SliceStrToIntM(ThirdStr, 1, 12)
				if err != nil {
					return "", err
				}

			}

			intSlice, err := SliceStrToIntM(SecondStr, -2, 31)
			if err != nil {
				return "", err
			}

			nowDay := d.Day()
			smallestDay := 0
			closestDay := 0
			firstNegativeDay := 0
			secondNegativeDay := 0

			//достаю из слайса необходимые дни
			for _, day := range intSlice {
				if smallestDay == 0 && day > 0 || day < smallestDay && day > 0 {
					smallestDay = day
				}
				if day > nowDay {
					if closestDay == 0 {
						closestDay = day
					}
					if closestDay != 0 && day-nowDay < closestDay-nowDay {
						closestDay = day
					}
				}
				if day == -1 {
					firstNegativeDay = day
				}
				if day == -2 {
					secondNegativeDay = day
				}
			}

			//иниц-ю получившиеся даты, от которых буду отталкиваться
			var closestDate time.Time

			if closestDay > 0 {
				closestDate = time.Date(d.Year(), d.Month(), closestDay, 0, 0, 0, 0, d.Location())

			}
			needDay := GetLastDayOfMonth(d.Year(), d.Month())
			firstNegativeDate := time.Date(d.Year(), d.Month(), needDay.Day(), 0, 0, 0, 0, d.Location())

			needDay = GetLastDayOfMonth(d.Year(), d.Month())
			secondNegativeDate := time.Date(d.Year(), d.Month(), needDay.Day()-1, 0, 0, 0, 0, d.Location())

			// 1. ищу даты в текущем месяце

			if closestDay > 0 {

				// 1.1 closestDay сущ-т в текущем месяце и негатива нет
				if closestDay == closestDate.Day() && firstNegativeDay+secondNegativeDay == 0 && closestDate.Format(f) > now.Format(f) && closestDate.Format(f) > d.Format(f) {
					// если третий параметр сущ-т && в нем есть текущий месяц || третьего параметра нет - дата подходит

					if shouldReturnDate(int(d.Month()), monthSlice) {
						return closestDate.Format(f), nil
					}
				}
				// 1.2 closestDay сущ-т в этом месяце и есть secondNegativeDay (-2)
				if closestDay == closestDate.Day() && secondNegativeDay < 0 && secondNegativeDate.Format(f) > now.Format(f) && secondNegativeDate.Format(f) > d.Format(f) {

					if secondNegativeDate.Format(f) > closestDate.Format(f) && secondNegativeDate.Day() > nowDay {
						if shouldReturnDate(int(d.Month()), monthSlice) {
							return secondNegativeDate.Format(f), nil
						}
					}
					if secondNegativeDate.Format(f) > closestDate.Format(f) && closestDate.Day() > nowDay {
						if shouldReturnDate(int(d.Month()), monthSlice) {
							return closestDate.Format(f), nil
						}
					}

				}
				// 1.3 closestDay сущ-т в этом месяце и есть firstNegativeDay (-1)
				if closestDay == closestDate.Day() && firstNegativeDay < 0 && firstNegativeDate.Format(f) > now.Format(f) && firstNegativeDate.Format(f) > d.Format(f) {

					if firstNegativeDate.Format(f) < closestDate.Format(f) && firstNegativeDate.Day() > nowDay {
						if shouldReturnDate(int(d.Month()), monthSlice) {
							return firstNegativeDate.Format(f), nil
						}
					}
					if firstNegativeDate.Format(f) > closestDate.Format(f) && closestDate.Day() > nowDay {
						if shouldReturnDate(int(d.Month()), monthSlice) {
							return closestDate.Format(f), nil
						}
					}

				}
			}

			// 1.4 closestDay не сущ-т в этом месяце, а негатив есть
			if firstNegativeDay+secondNegativeDay != 0 {

				//второй негатив подходит
				if secondNegativeDay < 0 && secondNegativeDate.Format(f) > now.Format(f) && secondNegativeDate.Format(f) > d.Format(f) {
					if shouldReturnDate(int(d.Month()), monthSlice) {
						return secondNegativeDate.Format(f), nil
					}
				}
				//первый негатив подходит
				if firstNegativeDay < 0 && firstNegativeDate.Format(f) > now.Format(f) && firstNegativeDate.Format(f) > d.Format(f) {
					if shouldReturnDate(int(d.Month()), monthSlice) {
						return firstNegativeDate.Format(f), nil
					}
				}
			}

			// 2. текущий месяц не подошел, смотрю следующие.

			//при использовании AddDate число и месяц могут сбиться, если числа нет в след-м месяце
			//поэтому AddDate использую только в getLastDayOfMonth()

			// 2.1 есть только негатив
			if smallestDay == 0 && firstNegativeDay+secondNegativeDay == 0 {

				//второй негатив cущ-т
				if secondNegativeDay < 0 {
					for {
						//перехожу на след. месяц
						needDay = GetLastDayOfMonth(secondNegativeDate.Year(), secondNegativeDate.Month()+1)
						secondNegativeDate = time.Date(secondNegativeDate.Year(), secondNegativeDate.Month()+1, needDay.Day()-1, 0, 0, 0, 0, secondNegativeDate.Location())
						//чекаю получившуюся дату
						if secondNegativeDate.Format(f) > now.Format(f) && secondNegativeDate.Format(f) > d.Format(f) {
							if shouldReturnDate(int(secondNegativeDate.Month()), monthSlice) {
								return secondNegativeDate.Format(f), nil
							}
						}
					}
				}
				//первый негатив cущ-т
				if firstNegativeDay < 0 {
					for {
						//перехожу на след. месяц
						needDay = GetLastDayOfMonth(firstNegativeDate.Year(), firstNegativeDate.Month()+1)
						firstNegativeDate = time.Date(firstNegativeDate.Year(), firstNegativeDate.Month()+1, needDay.Day(), 0, 0, 0, 0, secondNegativeDate.Location())
						//чекаю получившуюся дату
						if firstNegativeDate.Format(f) > now.Format(f) && firstNegativeDate.Format(f) > d.Format(f) {
							if shouldReturnDate(int(firstNegativeDate.Month()), monthSlice) {
								return firstNegativeDate.Format(f), nil
							}
						}
					}
				}
			}

			// 2.2 smallestDay сущ-т
			//нахожу элемент следующего месяца для будущего сравнения с датой
			lastDayOfMonth := GetLastDayOfMonth(d.Year(), d.Month()+1)
			//нахожу дату в след.месяце
			smallestDate := time.Date(lastDayOfMonth.Year(), lastDayOfMonth.Month(), smallestDay, 0, 0, 0, 0, lastDayOfMonth.Location())
			for {
				// если  значение smallestDay не поменялось при добавлении в дату - smallestDay существует в рассматр-м месяце
				if smallestDay == smallestDate.Day() {
					if smallestDate.Format(f) > now.Format(f) {
						if shouldReturnDate(int(smallestDate.Month()), monthSlice) {
							return smallestDate.Format(f), nil
						}
					}
				}
				lastDayOfMonth = GetLastDayOfMonth(lastDayOfMonth.Year(), lastDayOfMonth.Month()+1)
				smallestDate = time.Date(lastDayOfMonth.Year(), lastDayOfMonth.Month(), smallestDay, 0, 0, 0, 0, lastDayOfMonth.Location())
			}
		}
		return "", fmt.Errorf("что-то пошло не так")
	}
	return "", fmt.Errorf("что-то пошло не так2")
}
