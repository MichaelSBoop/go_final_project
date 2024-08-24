package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Scheduler() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		return db, fmt.Errorf("incorrect database path")
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbPath)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return db, err
	}
	if install {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
	}
	if err != nil {
		return db, err
	}
	return db, nil
}
