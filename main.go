package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	// "github.com/MichaelSBoop/go_final_project/handlers"

	_ "modernc.org/sqlite"

	"github.com/MichaelSBoop/go_final_project/handlers"
	"github.com/MichaelSBoop/go_final_project/storage"
	"github.com/joho/godotenv"
)

const Format string = "20060102"

var webDir string = "./web"

// Загружаем переменные окружения сразу при запуске main
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

// Проверка наличия и правильности порта
func checkPort() string {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		fmt.Println("no port specified in .env file, switching to port 7540")
		return ":" + port
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("incorrect port type")
		return ""
	}
	return "127.0.0.1:" + port
}

func main() {
	storage, err := storage.CreateStorage(Scheduler())
	if err != nil {
		fmt.Printf("database setup error:%v\n", err)
	}
	defer storage.DB.Close()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.HandleNextDate)
	// http.HandleFunc("/api/task", handlers.HandleTask)
	// http.HandleFunc("/api/tasks", handlers.HandleTasks)
	// http.HandleFunc("/api/task/done", handlers.HandleTaskDone)

	// Прослушиваем порт, стандартный или взятый из окружения
	if err := http.ListenAndServe(checkPort(), nil); err != nil {
		fmt.Printf("server setup error:%v\n", err)
	}
}
