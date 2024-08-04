package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func main() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	if err := http.ListenAndServe(checkPort(), nil); err != nil {
		fmt.Printf("server setup error:%s", err)
	}
}
