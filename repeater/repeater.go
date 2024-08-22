package repeater

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Константа стандартного временного формата
const Format string = "20060102"

// NextDate вычисляет дедлайн исполнения следующей задачи, опираясь на текущую дату, установленный пользователем дедлайн
// и правило повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	var response string
	repeatValues := strings.Split(repeat, " ")
	dateParsed, err := time.Parse(Format, date)
	if err != nil {
		return response, err
	}

	// Определяем поведение согласно полученному правилу
	if repeatValues[0] == "y" {
		nextYear := dateParsed.AddDate(1, 0, 0)
		for nextYear.Before(now) {
			nextYear = nextYear.AddDate(1, 0, 0)
		}
		return nextYear.Format(Format), nil
	}
	if repeatValues[0] == "d" && len(repeatValues) > 1 {
		days, err := strconv.Atoi(repeatValues[1])
		if days > 400 || days <= 0 {
			return response, fmt.Errorf("incorrect repetition format")
		}
		if err != nil {
			return response, err
		}
		var nextDay time.Time
		nextDay = dateParsed.AddDate(0, 0, days)
		for nextDay.Before(now) {
			nextDay = nextDay.AddDate(0, 0, days)
		}
		return nextDay.Format(Format), nil
	}
	if repeatValues[0] == "w" && len(repeatValues[0]) != 0 {
		var nextDay time.Time
		days := strings.Split(repeatValues[0], ",")
		for _, day := range days {
			dayInt, err := strconv.Atoi(day)
			if dayInt <= 0 || dayInt > 7 {
				return response, fmt.Errorf("incorrect repetition format")
			}
			if err != nil {
				return response, err
			}
			weekDay := int(now.Weekday())
			if weekDay < dayInt {
				nextDay = now.AddDate(0, 0, dayInt-weekDay)
				return nextDay.Format(Format), nil
			}
			nextDay = now
			for byte(nextDay.Weekday()) != day[0] {
				nextDay = nextDay.AddDate(0, 0, 1)
			}
			return nextDay.Format(Format), nil
		}
	}
	if repeatValues[0] == "m" {

	}
	return response, fmt.Errorf("incorrect repetition format")
}
