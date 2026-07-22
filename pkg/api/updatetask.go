package api

import (
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"net/http"
)

func updateTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task
	err := json.NewDecoder(req.Body).Decode(&task)
	defer req.Body.Close()
	if err != nil {
		err = fmt.Errorf("incorrect JSON. %w", err)
		writeJson(res, err)
		return
	}
	if task.Title == "" {
		err = fmt.Errorf("incorrect title. %w", err)
		writeJson(res, err)
		return
	}
	err = checkDate(&task)
	if err != nil {
		err = fmt.Errorf("incorrect date. %w", err)
		writeJson(res, err)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		err = fmt.Errorf("db error. %w", err)
		writeJson(res, err)
		return
	}
	writeJson(res, task)
}
