package handlers

import (
	"fmt"
	"net/http"

	jsr "github.com/MichaelSBoop/go_final_project/JSONResponse"
	"github.com/MichaelSBoop/go_final_project/storage"
)

const limit = 50

// HandleTasks формирует список задач на основе заданного лимита
func HandleTasks(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем правильность метода
		if r.Method != http.MethodGet {
			jsr.ErrorJSON(w, fmt.Errorf("incorrect method"), http.StatusBadRequest)
			return
		}
		// Получаем список задач
		tasks, err := s.GetTasks(limit)
		if err != nil {
			jsr.ErrorJSON(w, fmt.Errorf("failed to retrieve tasks from database: %v", err), http.StatusBadRequest)
			return
		}
		// Формируем JSON ответ
		jsonTasks := jsr.FormulateResponseTasks(tasks)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTasks)
		if err != nil {
			fmt.Println("failed to write data response")
		}
	}
}
