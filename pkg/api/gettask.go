package api

import (
	"main/pkg/db"
	"net/http"
)

func getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(res, err, 500)
		return
	}
	writeJson(res, task, 200)
}
