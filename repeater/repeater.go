package repeater

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет дедлайн исполнения следующей задачи, опираясь на текущую дату, установленный пользователем дедлайн
// и правило повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	var response string
	repeatValues := strings.Split(repeat, " ")
	dateParsed, err := time.Parse("20060102", date)
	if err != nil {
		return response, err
	}
	if repeatValues[0] == "y" {
		nextYear := dateParsed.AddDate(1, 0, 0)
		for nextYear.Before(now) {
			nextYear = nextYear.AddDate(1, 0, 0)
		}
		return nextYear.Format("20060102"), nil
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
		if days == 1 && now.Format("20060102") == date {
			nextDay = now
			return nextDay.Format("20060102"), nil
		}
		nextDay = dateParsed
		nextDay = nextDay.AddDate(0, 0, days)
		for nextDay.Before(now) || nextDay == now {
			nextDay = nextDay.AddDate(0, 0, days)
		}
		if nextDay.After(now) || nextDay == now {
			return nextDay.Format("20060102"), nil
		}
		return nextDay.Format("20060102"), nil
	}
	return response, fmt.Errorf("incorrect repetition format")
}
