package repeater

import (
	"fmt"
	"time"

	"github.com/MichaelSBoop/go_final_project/task"
)

// NewDatePost возвращает новую дату в соответствии с критериями POST-запроса:
// - если дата не указана, возвращается сегодняшняя дата
// - если дата указана позже сегодняшнего дня включительно, возвратить дату
// - если дата указана до сегодняшней и правило повторения отсутствует, возвращается сегодняшняя дата
// - если дата указана до сегодняшней и правило повторения присутствует, высчитываем новую дату после сегодняшнего дня включительно
func NewDatePost(task task.Task) (string, error) {
	var newDate string
	now := time.Now()
	if task.Date == "" {
		return now.Format(Format), nil
	}
	dateParsed, err := time.Parse(Format, task.Date)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}
	if dateParsed.Before(now) || task.Date != now.Format(Format) {
		if task.Repeat == "d 1" && task.Date == now.Format(Format) {
			return time.Now().Format(Format), nil
		}
		if task.Date == now.Format(Format) && task.Repeat != "" {
			return time.Now().Format(Format), nil
		}
		if task.Repeat != "" {
			newDate = task.Date
			if task.Repeat != "" {
				newDate, err = NextDate(now, newDate, task.Repeat)
				if err != nil {
					return "", err
				}
			}
			return newDate, nil
		} else {
			return time.Now().Format(Format), nil
		}
	} else if task.Date == now.Format(Format) && task.Repeat != "" {
		return time.Now().Format(Format), nil
	}
	return task.Date, nil
}
