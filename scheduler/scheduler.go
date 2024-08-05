package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Scheduler(db *sql.DB) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		fmt.Println("incorrect database path")
		return
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbPath)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install == true {
		db.Exec(`CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(64) NOT NULL DEFAULT '',
    comment VARCHAR(64) NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX scheduler_date ON scheduler (date);`)
	}
}
