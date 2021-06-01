package service

import (
	"errors"

	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(listId int, userId int, input todo.Item) (int, error) {
	_, err := s.listRepo.GetById(listId, userId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, input)
}

func (s *TodoItemService) GetAll(listId int, userId int) ([]todo.Item, error) {
	_, err := s.listRepo.GetById(listId, userId)
	if err != nil {
		return nil, errors.New("record not found")
	}
	return s.repo.GetAll(listId, userId)
}

func (s *TodoItemService) GetById(userId int, itemId int) (todo.Item, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId int, itemId int) error {
	if _, err := s.repo.GetById(userId, itemId); err != nil {
		return errors.New("record not found")
	}
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId int, itemId int, input todo.UpdateItemInput) error {
	if _, err := s.repo.GetById(userId, itemId); err != nil {
		return errors.New("record not found")
	}
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, itemId, input)
}
