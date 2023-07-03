package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusCompleted TodoStatus = "completed"
)

type Todo struct {
	Id          pgtype.UUID `json:"id"`
	Todo        string      `json:"todo"`
	Status      TodoStatus  `json:"status"`
	DateCreated time.Time   `json:"date_created"`
}

func connectToDatabase() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	pgURI := os.Getenv("PG_URI")
	if pgURI == "" {
		return nil, fmt.Errorf("PG_URI environment variable is not set: %w", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(pgURI)
	if err != nil {
		return nil, fmt.Errorf("problem parsing PG_URI: %w", err)
	}

	pgConnPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pgConnPool, nil
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	conn, err := connectToDatabase()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: cannot connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), "SELECT id, todo, status, date_created FROM todo")
	if err != nil {
		http.Error(w, "Internal server error: unable to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.Id, &todo.Todo, &todo.Status, &todo.DateCreated)
		if err != nil {
			http.Error(w, "Internal server error: error occurred during row iteration", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	conn, err := connectToDatabase()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: cannot connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	idInput := chi.URLParam(r, "id")
	id := pgtype.UUID{}
	err = id.Scan(idInput)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
		return
	}

	todo := Todo{}
	err = conn.QueryRow(context.Background(), "SELECT * FROM todo WHERE id = $1", id).Scan(&todo.Id, &todo.Todo, &todo.Status, &todo.DateCreated)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	conn, err := connectToDatabase()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: cannot connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	todo := Todo{}
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = conn.QueryRow(context.Background(), "INSERT INTO todo (todo) VALUES ($1) RETURNING id, status, date_created", todo.Todo).Scan(&todo.Id, &todo.Status, &todo.DateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	conn, err := connectToDatabase()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: cannot connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	idInput := chi.URLParam(r, "id")
	id := pgtype.UUID{}
	err = id.Scan(idInput)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
		return
	}

	todo := Todo{}
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = conn.Exec(context.Background(), "UPDATE todo SET todo = $1, status = $2 WHERE id = $3", todo.Todo, todo.Status, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = conn.QueryRow(context.Background(), "SELECT * FROM todo WHERE id = $1", id).Scan(&todo.Id, &todo.Todo, &todo.Status, &todo.DateCreated)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(todo)

}

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
