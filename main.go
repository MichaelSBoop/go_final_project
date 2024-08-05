package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	rep "github.com/MichaelSBoop/go_final_project/repeater"
	"github.com/joho/godotenv"
)

var webDir string = "./web"

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func checkPort() string {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		fmt.Println("no port specified in .env file")
		return ""
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("incorrect port type")
		return ""
	}
	return ":" + port
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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
	resp, err := rep.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

// func HandleTask(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "incorrect http method", http.StatusBadRequest)
// 		return
// 	}
// }

func main() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", HandleNextDate)
	//http.HandleFunc("/api/task", HandleTask)
	if err := http.ListenAndServe(checkPort(), nil); err != nil {
		fmt.Printf("server setup error:%s", err)
	}
}
