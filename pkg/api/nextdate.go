package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "20060102"

func ParseRepeat(repeat string) (string, int, error) {
	if strings.HasPrefix(repeat, "d ") {
		repeatArray := strings.Split(repeat, " ")
		if len(repeatArray) == 0 {
			return "", 0, fmt.Errorf("некорректный ввод")
		}
		letter := repeatArray[0]
		days, err := strconv.Atoi(repeatArray[1])
		if err != nil {
			return "", 0, fmt.Errorf("ошибка")
		}
		return letter, days, nil
	}
	return repeat, 0, nil
}

func DateParse(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse(timeFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге времени date: %w", err)
	}

	taskDays, err := NextDate(now, date.Format(timeFormat), repeat)
	if err != nil {
		return "", fmt.Errorf("ошибка в функции NextDate: %w", err)
	}

	return taskDays, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	var futureDate time.Time

	t, err := time.Parse(timeFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге даты: %w", err)
	}

	letter, days, err := ParseRepeat(repeat)
	if err != nil {
		return "", err
	}

	if letter == "d" {
		if days == 0 {
			return "", fmt.Errorf("не указан интервал в днях")
		}

		if days > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал дней: %d", days)
		}

		futureDate = t.AddDate(0, 0, days)
		for futureDate.Before(now) {
			futureDate = futureDate.AddDate(0, 0, days)
		}
		return futureDate.Format(timeFormat), nil
	}
	if letter == "y" {
		futureDate = t.AddDate(1, 0, 0)
		for futureDate.Before(now) {
			futureDate = futureDate.AddDate(1, 0, 0)
		}
		return futureDate.Format(timeFormat), nil
	} else {
		return "", fmt.Errorf("недопустимый символ: %s", letter)
	}
}
