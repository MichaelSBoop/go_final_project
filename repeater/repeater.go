package repeater

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	var response string
	repeatValues := strings.Split(repeat, " ")
	dateParsed, err := time.Parse("20060102", date)
	if err != nil {
		return response, err
	}
	if repeatValues[0] == "y" {
		dateParsed = dateParsed.AddDate(1, 0, 0)
		for dateParsed.Before(now) || dateParsed == now {
			dateParsed = dateParsed.AddDate(1, 0, 0)
		}
		return dateParsed.Format("20060102"), nil
	}
	if repeatValues[0] == "d" && len(repeatValues) > 1 {
		days, err := strconv.Atoi(repeatValues[1])
		if days > 400 {
			return response, fmt.Errorf("incorrect repetition format")
		}
		if err != nil {
			return response, err
		}
		dateParsed = dateParsed.AddDate(0, 0, days)
		for dateParsed.Before(now) || dateParsed == now {
			dateParsed = dateParsed.AddDate(0, 0, days)
		}
		return dateParsed.Format("20060102"), nil
	}
	return response, fmt.Errorf("incorrect repetition format")
}
