package repository

import (
	"github.com/fr13n8/todo-app/structs"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user structs.SignUpInput) (int, error)
	GetUser(username string) (structs.User, error)
	CreateSession(input structs.Session) error
}

type TodoList interface {
	Create(userId int, list structs.List) (int, error)
	GetAll(userId int) ([]structs.List, error)
	GetById(listId int, userId int) (structs.List, error)
	Delete(listId int, userId int) error
	Update(listId int, userId int, input structs.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input structs.Item) (int, error)
	GetAll(listId int, userId int) ([]structs.Item, error)
	GetById(userId int, itemId int) (structs.Item, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input structs.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
