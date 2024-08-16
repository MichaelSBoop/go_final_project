// Пакет storage предоставляет инструменты для взаимодействия с базой данных на основе методов
// созданной структуры хранилища
package storage

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/MichaelSBoop/go_final_project/task"
	_ "modernc.org/sqlite"
)

// Константа лимита передаваемых задач
const limit = 50

// Создание хранилища задач
type Storage struct {
	DB *sql.DB
}

func CreateStorage(db *sql.DB, err error) (Storage, error) {
	return Storage{DB: db}, nil
}

// AddTask добавляет задачу в созданную базу и возвращает id
func (s Storage) AddTask(task task.Task) (string, error) {
	resp, err := s.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", fmt.Errorf("failed to add task:%v", err)
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve task id:%v", err)
	}
	// Последний идентификатор возвращается в формате int64, приводим его в строку и возвращаем
	return strconv.FormatInt(id, 10), nil
}

// GetTasks позволяет получить набор задач с заданным лимитом
func (s Storage) GetTasks() ([]task.Task, error) {
	var tasks []task.Task
	rows, err := s.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task task.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to put data into task type: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to loop through tasks: %v", err)
	}
	if tasks == nil {
		tasks = []task.Task{}
	}
	return tasks, nil
}

// GetTask ищет в базе задачу по заданному id и возвращает её
func (s Storage) GetTask(id int) (task.Task, error) {
	var task task.Task
	row := s.DB.QueryRow("SELECT date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("no rows were found: %v", err)
		}
		return task, err
	}
	return task, nil
}

// ChangeTask меняет данные переданной задачи в базе
func (s Storage) ChangeTask(task task.Task) error {
	res, err := s.DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("no rows were affected")
	}
	return nil
}

// DeleteTask удаляет задачу из базы по id
func (s Storage) DeleteTask(id int) error {
	res, err := s.DB.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("no rows were affected")
	}
	return nil
}
