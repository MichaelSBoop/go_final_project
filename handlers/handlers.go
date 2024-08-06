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

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}
	nowString := r.FormValue("now")
	nowTime, err := time.Parse("20060102", nowString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	resp, err := repeater.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	switch r.Method {
	case http.MethodPost:
		//jsonResponse := make(map[string]string)
		var task types.Task
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			// jsonResponse["error"] = "failed to read request: " + err.Error()
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
			// jsonResponse["error"] = err.Error()
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
			// jsonResponse["error"] = "failed to read unmarshal: " + err.Error()
			// json.NewEncoder(w).Encode(jsonResponse)
			// return
			// jsonResponse["error"] = err.Error()
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
		}
		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "title is required"}
			json.NewEncoder(w).Encode(jsonErr)
			return
			// jsonResponse["error"] = "title is required"
			// json.NewEncoder(w).Encode(jsonResponse)
			// return
			// jsonResponse["error"] = "title is required"
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
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
			// jsonResponse["error"] = "bad time format: " + err.Error()
			// json.NewEncoder(w).Encode(jsonResponse)
			// return
			// jsonResponse["error"] = err.Error()
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
		}
		if dateParsed.Before(time.Now()) {
			if task.Repeat != "" {
				task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					jsonErr := map[string]string{"error": err.Error()}
					json.NewEncoder(w).Encode(jsonErr)
					return
					// jsonResponse["error"] = "failed to create repetition: " + err.Error()
					// json.NewEncoder(w).Encode(jsonResponse)
					// return
					// jsonResponse["error"] = err.Error()
					// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
					// return
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
			// jsonResponse["error"] = "failed to write to database: " + err.Error()
			// json.NewEncoder(w).Encode(jsonResponse)
			// return
			// jsonResponse["error"] = err.Error()
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
		}
		//jsonResponse["id"] = taskId
		jsonId, err := json.Marshal(map[string]string{"id": taskId})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
			// jsonResponse["error"] = "failed to marshal: " + err.Error()
			// json.NewEncoder(w).Encode(jsonResponse)
			// return
			// jsonResponse["error"] = err.Error()
			// http.Error(w, fmt.Sprint(jsonResponse), http.StatusBadRequest)
			// return
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
		idInt, err := strconv.Atoi(task.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": "incorrect id"}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
		if idInt > 500 {
			w.WriteHeader(http.StatusBadRequest)
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
	if task.Repeat == "" {
		if err := actions.DeleteTask(task.ID, db); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(jsonErr)
			return
		}
	} else {
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
