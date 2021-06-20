package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fr13n8/todo-app/pkg/repository"
	"github.com/fr13n8/todo-app/structs"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user structs.SignUpInput) (int, error)
	GenerateToken(user structs.User) ([]string, error)
	SignInUser(username, password, userAgent string) ([]string, error)
	ParseToken(token string) (*jwt.StandardClaims, error)
	RefreshToken(token string) ([]string, error)
	CreateSession(input structs.Session) error
}

type TodoList interface {
	Create(userId int, list structs.List) (int, error)
	GetAll(userId int) ([]structs.List, error)
	GetById(listId int, userId int) (structs.List, error)
	Delete(listId int, userId int) error
	Update(listId int, userId int, list structs.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, userId int, input structs.Item) (int, error)
	GetAll(listId int, userId int) ([]structs.Item, error)
	GetById(userId int, itemId int) (structs.Item, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input structs.UpdateItemInput) error
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
