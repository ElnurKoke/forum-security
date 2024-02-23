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
	log.Println("Running a web server on https://localhost:8080")
	err := http.ListenAndServeTLS(":8080", "secure/server.crt", "secure/server.key", handlers.Mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
