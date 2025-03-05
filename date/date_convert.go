package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const FormatDate = "20060102"

// Convert - пересобирает дату в удобный формат
func Convert(search string) (string, error) {

	result := fmt.Sprint(search[6:] + search[3:5] + search[:2])

	_, err := Validation(result)
	if err != nil {
		fmt.Println("Преобразовать дату не удалось ", err)
		return result, err
	}
	return result, nil
}

// Validation - проверка корректности переданной даты
func Validation(date string) (time.Time, error) {

	var dateTime time.Time

	dateTime, err := time.Parse(FormatDate, date)
	if err != nil {
		fmt.Println("Строковые данные даты не корректны. ", err)
		return dateTime, err
	}

	return dateTime, nil
}

// RepeatValidation - проверка корректности переданного значениядля повторений
func RepeatValidation(repeat string) error {

	if repeat == "" {
		return nil
	}

	repeatData := strings.Split(repeat, " ")

	switch string(repeatData[0]) {
	// Ежегодно
	case "y":
		if len(repeatData) == 1 {
			return nil

		}
	// Через несколько дней
	case "d":
		if len(repeatData) < 2 {
			err2 := errors.New("не указан интервал для повторений в днях")
			return err2
		}

		if len(repeatData) > 2 {
			err3 := errors.New("не верно указан интервал для повторений в днях")
			return err3
		}

		daysCount, err := strconv.Atoi(repeatData[1])
		if err != nil {
			fmt.Println("Неверный формат дней для повторений. ", err)
			return err
		}

		if daysCount > 400 {
			err4 := errors.New("превышен интервал для повторений в днях")
			return err4
		}
	// По дням недели
	case "w":
		weekDays := strings.Split(repeat, " ")

		// Проверка на наличие дня недели
		if len(weekDays) < 2 {
			err5 := errors.New("не верный день недели")
			fmt.Println(err5)
			return err5
		}

		// Проверка на наличие одного дня недели
		if len(weekDays[1]) == 1 {
			dayNumber, err := strconv.Atoi(weekDays[1])

			if err != nil {
				fmt.Println("не верное значение дня недели", err)
				return err
			}

			if 0 >= dayNumber || dayNumber >= 8 {
				err6 := errors.New("не верный день недели")
				fmt.Println(err6)
				return err6
			}
			// Если дней не один
		} else {

			for _, value := range strings.Split(weekDays[1], ",") {
				day, err := strconv.Atoi(value)

				if err != nil {
					fmt.Println("не верное значение дня недели", err)
					return err
				}

				if 0 >= day || day >= 8 {
					err7 := errors.New("не верный день недели")
					fmt.Println(err7)
					return err7
				}
			}
		}
		// TODO Дописать проверку по месяцам

	default:
		err8 := errors.New("неизвестное значение для повторений")
		fmt.Println(err8)
		return err8
	}

	return nil
}
