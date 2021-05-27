package repository

import (
	"github.com/fr13n8/todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.List) (int, error)
	GetAll(userId int) ([]todo.List, error)
	GetById(listId int, userId int) (todo.List, error)
	Delete(listId int, userId int) error
	Update(listId int, userId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input todo.Item) (int, error)
	GetAll(listId int, userId int) ([]todo.Item, error)
	GetById(userId int, itemId int) (todo.Item, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input todo.UpdateItemInput) error
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
