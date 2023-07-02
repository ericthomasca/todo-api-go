package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
		log.Fatal("Error loading .env file")
	}
	
	pgURI := os.Getenv("PG_URI")
	if pgURI == "" {
		return nil, fmt.Errorf("PG_URI environment variable is not set")
	}

	pgxConfig, err := pgxpool.ParseConfig(pgURI)
	if err != nil {
		return nil, fmt.Errorf("problem parsing PG_URI: %v", err)
	}
	
	pgConnPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return pgConnPool, nil
}

func main() {
	// Connect to database
	pgConnPool, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pgConnPool.Close()

	// GET all todos from database
	todos, err := pgConnPool.Query(context.Background(), "SELECT id, todo, status, date_created FROM todo")
	if err != nil {
		log.Fatalf("Unable to execute query: %v\n", err)
	}
	defer todos.Close()

	// Print each todo
	fmt.Println("Todos:")
	for todos.Next() {
		var todo Todo
		err := todos.Scan(&todo.Id, &todo.Todo, &todo.Status, &todo.DateCreated)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan todo: %v\n", err)
			continue
		}

		todoIdValue, err := todo.Id.Value()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem with UUID")
		}

		fmt.Printf("ID: %s\tTodo: %s\tStatus: %s\tDate Created: %s\n", todoIdValue, todo.Todo, todo.Status, todo.DateCreated.String())

		// TODO CRUD API calls
	}
 
	os.Exit(0)
}
