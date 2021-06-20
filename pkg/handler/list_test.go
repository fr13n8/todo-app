package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/fr13n8/todo-app/pkg/service"
	mockservice "github.com/fr13n8/todo-app/pkg/service/mocks"
	"github.com/fr13n8/todo-app/structs"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createList(t *testing.T) {
	type input struct {
		userId int
		list   structs.List
	}

	type mockBehavior func(s *mockservice.MockTodoList, input input)

	testTable := []struct {
		name                 string
		input                input
		mockBehavior         mockBehavior
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				list: structs.List{
					Title:       "title",
					Description: "description",
				},
			},
			inputBody:            `{"title":"title","description":"description"}`,
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Create(input.userId, input.list).Return(1, nil)
			},
		},
		{
			name: "OK_WithoutDescription",
			input: input{
				userId: 1,
				list: structs.List{
					Title: "title",
				},
			},
			inputBody:            `{"title":"title"}`,
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Create(input.userId, input.list).Return(1, nil)
			},
		},
		{
			name: "No input",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
			mockBehavior:         func(r *mockservice.MockTodoList, input input) {},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				list: structs.List{
					Title:       "title",
					Description: "description",
				},
			},
			inputBody:            `{"title":"title","description":"description"}`,
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Create(input.userId, input.list).Return(0, errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mockservice.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.input)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllList(t *testing.T) {
	type input struct {
		userId int
	}

	type mockBehavior func(mockservice *mockservice.MockTodoList, input input)

	testTable := []struct {
		name                 string
		input                input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":[{"id":1,"title":"title","description":"description"},{"id":2,"title":"title2","description":"description2"}]}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetAll(input.userId).Return([]structs.List{
					{
						Id:          1,
						Title:       "title",
						Description: "description",
					},
					{
						Id:          2,
						Title:       "title2",
						Description: "description2",
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
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetAll(input.userId).Return(nil, errors.New("record not found"))
			},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetAll(input.userId).Return(nil, errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mockservice.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.input)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.getAllList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists/", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getListById(t *testing.T) {
	type input struct {
		listId int
		userId int
	}

	type mockBehavior func(mockservice *mockservice.MockTodoList, input input)

	testTable := []struct {
		name                 string
		input                input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":{"id":1,"title":"title","description":"description"}}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetById(input.listId, input.userId).Return(structs.List{
					Id:          1,
					Title:       "title",
					Description: "description",
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
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetById(input.listId, input.userId).Return((structs.List{}), errors.New("item not found"))
			},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().GetById(input.listId, input.userId).Return(structs.List{}, errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mockservice.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.input)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/lists/:id", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.getListById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%d", testCase.input.listId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
	type input struct {
		listId int
		userId int
	}

	type mockBehavior func(mockservice *mockservice.MockTodoList, input input)

	testTable := []struct {
		name                 string
		input                input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Delete(input.listId, input.userId).Return(nil)
			},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Delete(input.listId, input.userId).Return(errors.New("service failure"))
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"record not found"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Delete(input.listId, input.userId).Return(errors.New("record not found"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mockservice.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.input)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/api/lists/:id", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.deleteList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%d", testCase.input.listId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	type input struct {
		listId int
		userId int
		list   structs.UpdateListInput
	}

	type mockBehavior func(mockservice *mockservice.MockTodoList, input input)

	testTable := []struct {
		name                 string
		input                input
		mockBehavior         mockBehavior
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				listId: 1,
				list: structs.UpdateListInput{
					Title:       stringPointer("title"),
					Description: stringPointer("description"),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"title":"title","description":"description"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Update(input.listId, input.userId, input.list).Return(nil)
			},
		},
		{
			name: "Ok_WithoutTitle",
			input: input{
				userId: 1,
				listId: 1,
				list: structs.UpdateListInput{
					Description: stringPointer("description"),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"description":"description"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Update(input.listId, input.userId, input.list).Return(nil)
			},
		},
		{
			name: "Ok_WithoutDescription",
			input: input{
				userId: 1,
				listId: 1,
				list: structs.UpdateListInput{
					Title: stringPointer("title"),
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
			inputBody:            `{"title":"title"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Update(input.listId, input.userId, input.list).Return(nil)
			},
		},
		{
			name: "Not found",
			input: input{
				userId: 1,
				listId: 1,
				list: structs.UpdateListInput{
					Description: stringPointer("description"),
					Title:       stringPointer("title"),
				},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"not found"}`,
			inputBody:            `{"description":"description","title":"title"}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Update(input.listId, input.userId, input.list).Return(errors.New("not found"))
			},
		},
		{
			name: "No inputs",
			input: input{
				userId: 1,
				listId: 1,
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
			mockBehavior:         func(r *mockservice.MockTodoList, input input) {},
		},
		{
			name: "Service failure",
			input: input{
				userId: 1,
				listId: 1,
				list: structs.UpdateListInput{
					Title:       stringPointer("title"),
					Description: stringPointer("description"),
				},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
			inputBody:            `{"title":"title","description":"description","done":true}`,
			mockBehavior: func(r *mockservice.MockTodoList, input input) {
				r.EXPECT().Update(input.listId, input.userId, input.list).Return(errors.New("service failure"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mockservice.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.input)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/api/lists/:id/", func(c *gin.Context) {
				c.Set(userCtx, testCase.input.userId)
			}, handler.updateList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%d/", testCase.input.listId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
