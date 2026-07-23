package api

import (
	"main/pkg/db"
	"net/http"
)

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	err := db.DeleteTask(id)
	writeJson(res, err, 500)
}
