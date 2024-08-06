package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/MichaelSBoop/go_final_project/types"
)

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []types.Task
	query := `SELECT * FROM scheduler ORDER BY date LIMIT 50`
	db, err := sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task types.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []types.Task{}
	}
	jsonRes, err := json.Marshal(map[string][]types.Task{"tasks": tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
