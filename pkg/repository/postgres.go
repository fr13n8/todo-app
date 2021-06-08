package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	todoListsTable     = "todo_lists"
	usersListsTable    = "users_lists"
	todoItemsTable     = "todo_items"
	listsItemsTable    = "lists_items"
	usersSessionsTable = "users_sessions"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	address := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Open("postgres", address)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
