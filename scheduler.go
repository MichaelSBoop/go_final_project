package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Scheduler ищет путь к базе данных, открывает с ней соединение или создаёт её с нуля, если её нет
func Scheduler() (*sql.DB, error) {

	// Проверка наличия пути для существующей базы данных
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbPath, ok := os.LookupEnv("TODO_DBFILE")
	if !ok {
		dbPath = "scheduler.db"
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbPath)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Если базы данных нет, создаём свою в корне проекта и создаём таблицу scheduler на основе запроса
	// SQL-файла scheduler.sql
	if install {
		file, err := os.Create(dbFile)
		if err != nil {
			panic(err)
		}
		file.Close()
		query, err := os.ReadFile("scheduler.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to read sql query:%v", err)
		}
		if len(query) == 0 {
			return nil, fmt.Errorf("sql query is empty")
		}
		queryString := string(query)
		_, err = db.Exec(queryString)
		if err != nil {
			return nil, fmt.Errorf("failed to create scheduler table:%v", err)
		}
	}
	return db, nil
}
