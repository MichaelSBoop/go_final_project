package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MichaelSBoop/go_final_project/actions"
	"github.com/MichaelSBoop/go_final_project/repeater"
)

func HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	id := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := map[string]string{"error": "incorrect id"}
		json.NewEncoder(w).Encode(jsonErr)
		return
	}
	task, err := actions.GetTask(id, db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(jsonErr)
		return
	}
	if task.Repeat != "" {
		now := time.Now()
		newDate, err := repeater.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		task.Date = newDate
		if err = actions.ChangeTask(task, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
	} else {
		if err := actions.DeleteTask(task.ID, db); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
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
