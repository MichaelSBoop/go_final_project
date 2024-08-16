package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MichaelSBoop/go_final_project/encode"
	rep "github.com/MichaelSBoop/go_final_project/repeater"
	"github.com/MichaelSBoop/go_final_project/storage"
	"github.com/MichaelSBoop/go_final_project/task"
)

// HandleTask обрабатывает GET, POST, PUT и DELETE http-запросы, обращаясь к базе данных
func SingleTask(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// Добавление задачи
		case http.MethodPost:
			// Считываем тело в буфер и проводим ряд проверок на возможные ошибки
			var task task.Task
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to read body: %v", err), http.StatusBadRequest)
				return
			}
			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to unmarshal data: %v", err), http.StatusBadRequest)
				return
			}
			if task.Title == "" {
				encode.ErrorJSON(w, fmt.Errorf("title is required"), http.StatusBadRequest)
				return
			}
			// Записываем дату в задачу в зависимости от наличия правила повторения и самой даты
			newDate, err := rep.PostLogic(task)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to calculate new date: %v", err), http.StatusBadRequest)
				return
			}
			task.Date = newDate
			// Добавляем задачу в базу данных и возвращаем её id
			taskId, err := s.AddTask(task)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to add task: %v", err), http.StatusBadRequest)
				return
			}
			// Формулируем JSON для записи
			jsonId := encode.FormulateResponseID("id", taskId)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(jsonId)
			if err != nil {
				fmt.Println("failed to write data response")
			}
		// Получение задачи
		case http.MethodGet:
			// Считываем id
			var task task.Task
			id := r.URL.Query().Get("id")
			if id == "" {
				encode.ErrorJSON(w, fmt.Errorf("id is required"), http.StatusBadRequest)
				return
			}
			// Получаем задачу
			task, err := s.GetTask(id)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to retrieve task: %v", err), http.StatusBadRequest)
				return
			}
			// Формируем JSON для записи
			jsonTask := encode.FormulateResponseTask(task)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonTask)
			if err != nil {
				fmt.Println("failed to write data response")
			}
		// Обновление задачи
		case http.MethodPut:
			// Для обновления задачи считываем тело в буфер и проводим те же проверки, которые используем для записи
			var task task.Task
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				encode.ErrorJSON(w, err, http.StatusBadRequest)
				return
			}
			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to unmarshal data:%v", err), http.StatusBadRequest)
				return
			}
			if task.ID == "" {
				encode.ErrorJSON(w, fmt.Errorf("id is required"), http.StatusBadRequest)
				return
			}
			_, err = strconv.Atoi(task.ID)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("incorrect id: %v", err), http.StatusBadRequest)
				return
			}
			_, err = s.GetTask(task.ID)
			if err != nil {
				encode.ErrorJSON(w, fmt.Errorf("incorrect id: %v", err), http.StatusBadRequest)
				return
			}
			if task.Title == "" {
				encode.ErrorJSON(w, fmt.Errorf("title is required"), http.StatusBadRequest)
				return
			}
			if task.Date == "" {
				task.Date = time.Now().Format(rep.Format)
			}
			dateParsed, err := time.Parse(rep.Format, task.Date)
			if err != nil {
				encode.ErrorJSON(w, err, http.StatusBadRequest)
				return
			}
			if dateParsed.Before(time.Now()) {
				if task.Repeat != "" {
					task.Date, err = rep.NextDate(time.Now(), task.Date, task.Repeat)
					if err != nil {
						encode.ErrorJSON(w, fmt.Errorf("failed to set next date: %v", err), http.StatusBadRequest)
						return
					}
				} else {
					task.Date = time.Now().Format(rep.Format)
				}
			}
			// Изменяем задачу в базе
			if err = s.ChangeTask(task); err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to change task data: %v", err), http.StatusBadRequest)
				return
			}
			// Формируем JSON для записи
			jsonEmpty := encode.FormulateResponseEmpty()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonEmpty)
			if err != nil {
				fmt.Println("failed to write data response")
			}
		// Удаление задачи
		case http.MethodDelete:
			// Получаем id и и на его основе удаляем задачу из базы
			id := r.URL.Query().Get("id")
			if id == "" {
				encode.ErrorJSON(w, fmt.Errorf("id is required"), http.StatusBadRequest)
				return
			}
			if err := s.DeleteTask(id); err != nil {
				encode.ErrorJSON(w, fmt.Errorf("failed to delete task: %v", err), http.StatusBadRequest)
				return
			}
			// Формируем JSON ответ
			jsonEmpty := encode.FormulateResponseEmpty()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(jsonEmpty)
			if err != nil {
				fmt.Println("failed to write data response")
			}
		}
	}
}
