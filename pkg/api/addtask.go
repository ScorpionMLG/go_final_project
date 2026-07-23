package api

import (
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"net/http"
	"time"
)

func addTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task
	err := json.NewDecoder(req.Body).Decode(&task)
	defer req.Body.Close()
	if err != nil {
		err = fmt.Errorf("incorrect JSON. %w", err)
		writeJson(res, err, 400)
		return
	}
	if task.Title == "" {
		err = fmt.Errorf("incorrect title. %w", err)
		writeJson(res, err, 400)
		return
	}
	err = checkDate(&task)
	if err != nil {
		err = fmt.Errorf("incorrect date. %w", err)
		writeJson(res, err, 400)
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		err = fmt.Errorf("db error. %w", err)
		writeJson(res, err, 500)
		return
	}
	task.ID = fmt.Sprint(id)
	writeJson(res, task, 201)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowFormated := now.Format(DateFormat)
	if task.Date == "" {
		task.Date = nowFormated
	}
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("bad date. %w", err)
	}
	nowDay, err := time.Parse(DateFormat, nowFormated)
	if err != nil {
		return fmt.Errorf("bad date. %w", err)
	}
	if nowDay.After(t) {
		if len(task.Repeat) != 0 {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("incorrect repeat. %w", err)
			}
			task.Date = next
		}
		if len(task.Repeat) == 0 {
			task.Date = nowFormated
		}
	}
	return nil
}
