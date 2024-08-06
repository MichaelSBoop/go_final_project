package actions

import (
	"database/sql"
	"strconv"
)

func DeleteTask(id string, db *sql.DB) error {
	_, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return err
	}
	return nil
}
