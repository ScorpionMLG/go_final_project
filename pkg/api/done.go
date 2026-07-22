package api

import (
	"main/pkg/db"
	"net/http"
	"time"
)

func doneTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(res, err)
		return
	}
	if task.Repeat == "" {
		err := db.DeleteTask(id)
		writeJson(res, err)
		return
	}
	nDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(res, err)
		return
	}
	task.Date = nDate
	err = db.UpdateTaskDate(task)
	writeJson(res, err)
}
