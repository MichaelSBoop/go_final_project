package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MichaelSBoop/go_final_project/handlers"
	"github.com/MichaelSBoop/go_final_project/scheduler"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
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

func main() {
	db, err := scheduler.Scheduler()
	if err != nil {
		fmt.Printf("database setup error:%v\n", err)
	}
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.HandleNextDate)
	http.HandleFunc("/api/task", handlers.HandleTask)
	http.HandleFunc("/api/tasks", handlers.HandleTasks)
	http.HandleFunc("/api/task/done", handlers.HandleTaskDone)
	if err := http.ListenAndServe(checkPort(), nil); err != nil {
		fmt.Printf("server setup error:%v\n", err)
	}
}
