package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MichaelSBoop/go_final_project/actions"
	"github.com/MichaelSBoop/go_final_project/repeater"
	"github.com/MichaelSBoop/go_final_project/types"
)

func HandleTask(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	switch r.Method {
	case http.MethodPost:
		var task types.Task
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "title is required"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}
		dateParsed, err := time.Parse("20060102", task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if dateParsed.Before(time.Now()) {
			if task.Repeat != "" {
				task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					jsonErr := map[string]string{"error": err.Error()}
					json.NewEncoder(w).Encode(jsonErr)
					return
				}
			} else {
				task.Date = time.Now().Format("20060102")
			}
		}
		taskId, err := actions.AddTask(task, db)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		jsonId, err := json.Marshal(map[string]string{"id": taskId})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonId)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		task, err := actions.GetTask(id, db)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		jsonTask, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonTask)

	case http.MethodPut:
		var task types.Task
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if task.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "id is required"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		_, err = strconv.Atoi(task.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "incorrect id"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		_, err = actions.GetTask(task.ID, db)
		if err != nil {
			jsonErr := map[string]string{"error": "incorrect id"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "title is required"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}
		dateParsed, err := time.Parse("20060102", task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if dateParsed.Before(time.Now()) {
			if task.Repeat != "" {
				task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					jsonErr := map[string]string{"error": err.Error()}
					json.NewEncoder(w).Encode(jsonErr)
					return
				}
			} else {
				task.Date = time.Now().Format("20060102")
			}
		}
		if err = actions.ChangeTask(task, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		jsonResp, err := json.Marshal(struct{}{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if _, err := strconv.Atoi(id); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "id is required"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if err := actions.DeleteTask(id, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		jsonResp, err := json.Marshal(struct{}{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
