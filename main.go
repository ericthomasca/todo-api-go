package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/todos", getTodos)
	r.Get("/todos/{id},", getTodo)
	r.Post("/todos", createTodo)
	r.Put("/todos/{id}", updateTodo)

	fmt.Println("Serving on http://localhost:3420...")
	err := http.ListenAndServe(":3420", r)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
