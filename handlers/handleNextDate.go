// Пакет handlers реализует обработчики событий
package handlers

import (
	"fmt"
	"net/http"
	"time"

	rep "github.com/MichaelSBoop/go_final_project/repeater"
)

// HandleNextDate вычисляет и передаёт дату следующего дедлайна для задачи в виде ответа на GET-запрос
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	// Проверяем правильность метода
	if r.Method != http.MethodGet {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}
	// Читаем значения из запроса
	nowString := r.FormValue("now")
	nowTime, err := time.Parse(rep.Format, nowString)
	if err != nil {
		http.Error(w, "failed to parse time:"+err.Error(), http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	// Создаём новую дату и записываем её в ответ
	resp, err := rep.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, "failed to calculate next date:"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("failed to write response")
	}
}
