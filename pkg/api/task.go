package api

import "net/http"

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodPut:
		updateTaskHandler(res, req)
	case http.MethodDelete:
		deleteTaskHandler(res, req)
	default:
		http.Error(res, "incorrect request: "+req.Method, http.StatusMethodNotAllowed)
	}
}
