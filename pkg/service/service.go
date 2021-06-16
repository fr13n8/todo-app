package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todo.SignUpInput) (int, error)
	GenerateToken(user todo.User) ([]string, error)
	SignInUser(username, password, userAgent string) ([]string, error)
	ParseToken(token string) (*jwt.StandardClaims, error)
	RefreshToken(token string) ([]string, error)
	CreateSession(input todo.Session) error
}

type TodoList interface {
	Create(userId int, list todo.List) (int, error)
	GetAll(userId int) ([]todo.List, error)
	GetById(listId int, userId int) (todo.List, error)
	Delete(listId int, userId int) error
	Update(listId int, userId int, list todo.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, userId int, input todo.Item) (int, error)
	GetAll(listId int, userId int) ([]todo.Item, error)
	GetById(userId int, itemId int) (todo.Item, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input todo.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
