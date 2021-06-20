package service

import (
	"errors"

	"github.com/fr13n8/todo-app/pkg/repository"
	"github.com/fr13n8/todo-app/structs"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list structs.List) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]structs.List, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(listId int, userId int) (structs.List, error) {
	return s.repo.GetById(listId, userId)
}

func (s *TodoListService) Delete(listId int, userId int) error {
	if _, err := s.repo.GetById(listId, userId); err != nil {
		return errors.New("record not found")
	}
	return s.repo.Delete(listId, userId)
}

func (s *TodoListService) Update(listId int, userId int, input structs.UpdateListInput) error {
	if _, err := s.repo.GetById(listId, userId); err != nil {
		return errors.New("record not found")
	}
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(listId, userId, input)
}
