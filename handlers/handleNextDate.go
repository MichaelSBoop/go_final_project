package handlers

import (
	"net/http"
	"time"

	"github.com/MichaelSBoop/go_final_project/repeater"
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
