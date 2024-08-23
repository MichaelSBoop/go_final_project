package encode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MichaelSBoop/go_final_project/task"
)

// TODO: сделать интерфейс, формирующий запросы на основе разных типов

// Функции FormulateResponse формируют значение для отправки и осуществляет проверку сериализации в JSON

func FormulateResponseID(data string) []byte {
	resMap := make(map[string]string)
	resMap["id"] = data
	jsonRes, err := json.Marshal(resMap)
	if err != nil {
		fmt.Printf("failed to marshal data: %v", err)
		return []byte{}
	}
	return jsonRes
}

func FormulateResponseTasks(tasks []task.Task) []byte {
	resMap := make(map[string][]task.Task)
	resMap["tasks"] = tasks
	jsonRes, err := json.Marshal(resMap)
	if err != nil {
		fmt.Printf("failed to marshal data: %v", err)
		return []byte{}
	}
	return jsonRes
}

func FormulateResponseTask(task task.Task) []byte {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		fmt.Printf("failed to marshal data: %v", err)
		return []byte{}
	}
	return jsonTask
}

func FormulateResponseEmpty() []byte {
	jsonEmpty, err := json.Marshal(struct{}{})
	if err != nil {
		fmt.Printf("failed to marshal data: %v", err)
		return []byte{}
	}
	return jsonEmpty
}

// ErrorJSON переводит ошибки в формат JSON и записывает их в ответ;
// аналогично http.Error()
func ErrorJSON(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	jsonErr := map[string]string{"error": err.Error()}
	json.NewEncoder(w).Encode(jsonErr)
}
