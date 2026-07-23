package api

import (
	"main/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.Tasks(50, req.FormValue("search"))
	if err != nil {
		writeJson(res, err, 500)
		return
	}
	writeJson(res, TasksResp{
		Tasks: tasks,
	}, 200)
}
