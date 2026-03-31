package main

import (
	"log"
	"net/http"

	"vdisplay/internal/api"
	"vdisplay/internal/service"
)

func main() {
	manager := service.NewSessionManager()
	handler := &api.Handler{Manager: manager}

	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	log.Println("Starting Virtual Display API on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
