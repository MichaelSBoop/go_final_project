package actions

import (
	"database/sql"
	"strconv"

	"github.com/MichaelSBoop/go_final_project/types"
)

func GetTask(id string, db *sql.DB) (types.Task, error) {
	var task types.Task
	_, err := strconv.Atoi(id)
	if err != nil {
		return types.Task{}, err
	}
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return types.Task{}, err
	}
	return task, nil
}
