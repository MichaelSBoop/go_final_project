package storage

import "database/sql"

// Создание хранилища задач

type Storage struct {
	DB *sql.DB
}

func CreateStorage(db *sql.DB, err error) (Storage, error) {
	return Storage{DB: db}, nil
}

//func (s Storage) AddTask()
