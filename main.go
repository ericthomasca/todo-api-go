package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
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
	
	pgURI := os.Getenv("PG_URI");
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

func getTodosFromRows() ([]Todo, error) {
	// Connect to database
	pgConnPool, err := connectToDatabase()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	defer pgConnPool.Close()

	// GET all todos from database
	rows, err := pgConnPool.Query(context.Background(), "SELECT id, todo, status, date_created FROM todo")
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}
	defer rows.Close()
	
	var todos []Todo

	for rows.Next() {
		var todo Todo
		
		err := rows.Scan(&todo.Id, &todo.Todo, &todo.Status, &todo.DateCreated)
		if err != nil {
			return nil, fmt.Errorf("error occurred during row iteration: %w", err)	
		}
		
		// TODO use this code later
		// todoIdValue, err := todo.Id.Value()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Problem with UUID")
		// }

		todos = append(todos, todo)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	if len(todos) == 0 {
		return nil, fmt.Errorf("no todos found: %w", err)
	}

	return todos, nil
}

func main() {
	// router := chi.NewRouter()
	// router.Use(middleware.Logger)
	
	// router.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })

	// http.ListenAndServe(":3420", router)

	todos, err := getTodosFromRows()
	if err != nil {
		log.Printf("Error retrieving todos: %v", err)
	}

	fmt.Println(todos[0].Todo)
	fmt.Println(todos[1].Todo)
	fmt.Println(todos[2].Todo)
}
