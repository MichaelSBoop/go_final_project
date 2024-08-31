package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/MichaelSBoop/go_final_project/handlers"
	"github.com/MichaelSBoop/go_final_project/storage"
	"github.com/joho/godotenv"
)

var webDir string = "./web"

// Загружаем переменные окружения сразу при запуске main
func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}
}

// Проверка наличия и правильности порта
func checkPort() string {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		fmt.Println("no port specified in .env file, switching to port 7540")
		return ":7540"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("incorrect port type, switching to 7540")
		return ":7540"
	}
	return ":" + port
}

func main() {
	// Создаём базу на основе данных scheduler: проверка на наличие базы и таблицы внутри
	db, err := storage.CreateStorage(Scheduler())
	if err != nil {
		fmt.Printf("database setup error:%v\n", err)
	}
	defer db.DB.Close()

	// TODO: добавить mux или chi для обработчиков
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDate)
	http.HandleFunc("/api/task", handlers.SingleTask(db))
	http.HandleFunc("/api/tasks", handlers.MultipleTasks(db))
	http.HandleFunc("/api/task/done", handlers.TaskDone(db))

	// Прослушиваем порт, стандартный или взятый из окружения
	port := checkPort()
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("server setup error:%v\n", err)
	}
}
