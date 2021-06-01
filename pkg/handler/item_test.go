package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/service"
	mockservice "github.com/fr13n8/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createItem(t *testing.T) {
	type input struct {
		userId int
		listId int
		item   todo.Item
	}

	type mockBehavior func(r *mockservice.MockTodoItem, input input)

	testTable := []struct {
		name                 string
		item                 input
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
		mockBehavior         mockBehavior
	}{
		{
			name: "Ok",
			item: input{
				item: todo.Item{
					Title:       "title",
					Description: "description",
				},
				userId: 1,
				listId: 1,
			},
			inputBody:            `{"title": "title", "description": "description"}`,
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Create(input.listId, input.userId, input.item).Return(1, nil)
			},
		},
		{
			name: "No inputs",
			item: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
			mockBehavior:         func(r *mockservice.MockTodoItem, input input) {},
		},
		{
			name: "Service failure",
			item: input{
				item: todo.Item{
					Title:       "title",
					Description: "description",
				},
				userId: 1,
				listId: 1,
			},
			inputBody:            `{"title": "title", "description": "description"}`,
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Create(input.listId, input.userId, input.item).Return(0, errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mockservice.NewMockTodoItem(c)
			testCase.mockBehavior(item, testCase.item)

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/api/:id/items", func(c *gin.Context) {
				c.Set(userCtx, testCase.item.userId)
			}, handler.createItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/%d/items", testCase.item.listId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllItems(t *testing.T) {
	type input struct {
		userId int
		listId int
	}

	type mockBehavior func(r *mockservice.MockTodoItem, input input)

	testTable := []struct {
		name                 string
		input                input
		expectedStatusCode   int
		expectedResponseBody string
		mockBehavior         mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":[{"id":1,"title":"title","description":"description","done":false},{"id":2,"title":"title2","description":"description2","done":true}]}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetAll(input.listId, input.userId).Return([]todo.Item{
					{
						Id:          1,
						Title:       "title",
						Description: "description",
						Done:        false,
					},
					{
						Id:          2,
						Title:       "title2",
						Description: "description2",
						Done:        true,
					},
				}, nil)
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"record not found"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetAll(input.listId, input.userId).Return(nil, errors.New("record not found"))
			},
		},
		{
			name: "Service failure",
			input: input{
				listId: 1,
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetAll(input.listId, input.userId).Return(nil, errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mockservice.NewMockTodoItem(c)
			testCase.mockBehavior(item, testCase.input)

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/lists/:id/items/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.getAllItems)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%d/items/", testCase.input.listId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getItemById(t *testing.T) {
	type input struct {
		itemId int
		userId int
	}

	type mockBehavior func(r *mockservice.MockTodoItem, input input)

	testTable := []struct {
		name                 string
		input                input
		expectedStatusCode   int
		expectedResponseBody string
		mockBehavior         mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				itemId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":{"id":1,"title":"title","description":"description","done":false}}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetById(input.userId, input.itemId).Return(todo.Item{
					Id:          1,
					Title:       "title",
					Description: "description",
					Done:        false,
				}, nil)
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"item not found"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetById(input.userId, input.itemId).Return((todo.Item{}), errors.New("item not found"))
			},
		},
		{
			name: "Service failure",
			input: input{
				itemId: 1,
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().GetById(input.userId, input.itemId).Return((todo.Item{}), errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mockservice.NewMockTodoItem(c)
			testCase.mockBehavior(item, testCase.input)

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/items/:id/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.getItemById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/items/%d/", testCase.input.itemId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	type input struct {
		userId int
		itemId int
	}

	type mockBehavior func(r *mockservice.MockTodoItem, input input)

	testTable := []struct {
		name                 string
		input                input
		expectedStatusCode   int
		expectedResponseBody string
		mockBehavior         mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				itemId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Delete(input.userId, input.itemId).Return(nil)
			},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				itemId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Delete(input.userId, input.itemId).Return(errors.New("service failure"))
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"record not found"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Delete(input.userId, input.itemId).Return(errors.New("record not found"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mockservice.NewMockTodoItem(c)
			testCase.mockBehavior(item, testCase.input)

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/api/items/:id/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.deleteItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/items/%d/", testCase.input.itemId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateItem(t *testing.T) {
	type input struct {
		userId int
		itemId int
		item   todo.UpdateItemInput
	}

	type mockBehavior func(s *mockservice.MockTodoItem, input input)

	testTable := []struct {
		name                 string
		input                input
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Title:       stringPointer("title"),
					Description: stringPointer("description"),
					Done:        boolPointer(true),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"title":"title","description":"description","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(nil)
			},
		},
		{
			name: "Ok_WithoutTitle",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Description: stringPointer("description"),
					Done:        boolPointer(true),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"description":"description","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(nil)
			},
		},
		{
			name: "Ok_WithoutDescription",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Title: stringPointer("title"),
					Done:  boolPointer(true),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"title":"title","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(nil)
			},
		},
		{
			name: "Ok_WithoutDone",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Description: stringPointer("description"),
					Title:       stringPointer("title"),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"description":"description","title":"title"}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(nil)
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Description: stringPointer("description"),
					Title:       stringPointer("title"),
					Done:        boolPointer(true),
				},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"not found"}`,
			inputBody:            `{"description":"description","title":"title","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(errors.New("not found"))
			},
		},
		{
			name: "No inputs",
			input: input{
				userId: 1,
				itemId: 1,
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
			mockBehavior:         func(r *mockservice.MockTodoItem, input input) {},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				itemId: 1,
				item: todo.UpdateItemInput{
					Title:       stringPointer("title"),
					Description: stringPointer("description"),
					Done:        boolPointer(true),
				},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			inputBody:            `{"title":"title","description":"description","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoItem, input input) {
				r.EXPECT().Update(input.userId, input.itemId, input.item).Return(errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mockservice.NewMockTodoItem(c)
			testCase.mockBehavior(item, testCase.input)

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/api/items/:id/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.updateItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/items/%d/", testCase.input.itemId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
