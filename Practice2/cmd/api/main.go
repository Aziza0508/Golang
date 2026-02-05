package main

import (
	"log"
	"net/http"
	"practice2/internal/handlers"
	"practice2/internal/middleware"
	"practice2/internal/storage"
)

func main() {
	taskStorage := storage.NewTaskStorage()
	taskHandler := handlers.NewTaskHandler(taskStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", taskHandler.HandleTasks)

	handler := middleware.Logging(middleware.Auth(mux))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
