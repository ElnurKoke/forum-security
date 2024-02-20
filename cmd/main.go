package main

import (
	"forum/internal/handler"
	"forum/internal/service"
	"forum/internal/storage"
	"log"
	"net/http"
)

func main() {
	db := storage.InitDB()
	store := storage.NewStorage(db)
	services := service.NewService(store)
	handlers := handler.NewHandler(services)
	handlers.InitRoutes()
	log.Println("Running a web server on http://localhost:8080")
	http.ListenAndServe(":8080", handlers.Mux)
}
