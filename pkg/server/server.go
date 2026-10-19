package server

import (
	"log"
	"main/pkg/api"
	"net/http"
	"os"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	mux := http.NewServeMux()
	api.Init(mux)

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	log.Printf("Server starting on port %s...", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}
	return nil
}
