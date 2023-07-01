package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusCompleted TodoStatus = "completed"
)

type Todo struct {
	Id          pgxUUID.UUID `json:"id"`
	Todo        string       `json:"todo"`
	Status      TodoStatus   `json:"status"`
	DateCreated time.Time    `json:"date_created"`
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Parse PG_URI from .env
	pgURI := os.Getenv("PG_URI")
	if pgURI == "" {
		log.Fatal("PG_URI environment variable is not set")
	}

	// Parse config
	pgxConfig, err := pgxpool.ParseConfig(pgURI)
	if err != nil {
		log.Fatal("Problem parsing PG_URI")
	}

	// Set up pgxUUID data type
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	// Connect to database
	pgxConnPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatal("Unable to connect to database")
	}
	defer pgxConnPool.Close()

	// Get all todos from database
	todos, err := pgxConnPool.Query(context.Background(), "SELECT id, todo, status, date_created FROM todo")
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
		fmt.Printf("Todo: %s\tID: %T\tStatus: %s\tDate Created: %s\n",todo.Id, todo.Todo, todo.Status, todo.DateCreated.String())
		// TODO fix printing of UUID
	}

	// TODO CRUD API calls

	os.Exit(0)
}
