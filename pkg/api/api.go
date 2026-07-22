package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/signin", authHandler)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(doneTaskHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
}

func writeJson(res http.ResponseWriter, data any) {
	if data == nil {
		data = struct{}{}
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		data = fmt.Errorf("serialization error. %w", err)
	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err, isError := data.(error)
	if isError || err != nil {
		type errorJSON struct {
			Error string `json:"error"`
		}
		var errorStruct errorJSON
		errorStruct.Error = fmt.Sprint(data)
		marshal, _ = json.Marshal(errorStruct)
	}
	res.Write(marshal)
}
