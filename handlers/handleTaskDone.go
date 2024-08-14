package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jsr "github.com/MichaelSBoop/go_final_project/JSONResponse"
	rep "github.com/MichaelSBoop/go_final_project/repeater"
	"github.com/MichaelSBoop/go_final_project/storage"
	"github.com/MichaelSBoop/go_final_project/task"
)

// HandleTaskDone помечает задачу как выполненную;
// в зависимости от наличия повторений задача может быть обновлена или удалена
func HandleTaskDone(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем правильность http-метода
		if r.Method != http.MethodPost {
			jsr.ErrorJSON(w, fmt.Errorf("incorrect method"), http.StatusBadRequest)
			return
		}
		var task task.Task
		// Получаем задачу по id и проверяем на наличие ошибок
		id := r.URL.Query().Get("id")
		if id == "" {
			jsr.ErrorJSON(w, fmt.Errorf("id is required"), http.StatusBadRequest)
			return
		}
		_, err := strconv.Atoi(id)
		if err != nil {
			jsr.ErrorJSON(w, fmt.Errorf("incorrect id: %v", err), http.StatusBadRequest)
			return
		}
		task, err = s.GetTask(id)
		if err != nil {
			jsr.ErrorJSON(w, fmt.Errorf("failed to retrieve task: %v", err), http.StatusBadRequest)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format(rep.Format)
		}
		// Если в задаче указано правило повторения, вычисляется следующая дата выполнения;
		// иначе задача считается как одноразовая и удаляется
		if task.Repeat != "" {
			newDate, err := rep.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				jsr.ErrorJSON(w, err, http.StatusInternalServerError)
				return
			}
			task.Date = newDate
			if err = s.ChangeTask(task); err != nil {
				jsr.ErrorJSON(w, fmt.Errorf("failed to change task data: %v", err), http.StatusBadRequest)
				return
			}
		} else {
			if err := s.DeleteTask(id); err != nil {
				jsr.ErrorJSON(w, fmt.Errorf("failed to delete task: %v", err), http.StatusBadRequest)
				return
			}
		}
		// Формулируем JSON для ответа
		jsonEmpty := jsr.FormulateResponseEmpty()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonEmpty)
		if err != nil {
			fmt.Println("failed to write data response")
		}
	}
}
