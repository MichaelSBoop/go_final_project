package actions

import (
	"database/sql"
	"strconv"

	"github.com/MichaelSBoop/go_final_project/types"
)

func AddTask(task types.Task, db *sql.DB) (string, error) {
	resp, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}
