package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/service"
	mock_service "github.com/fr13n8/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "username": "test", "password": "test"}`,
			inputUser: todo.User{
				Name:     "Test",
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user todo.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"ID":1}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"username": "test", "password": "test"}`,
			mockBehavior:         func(r *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name": "Test", "username": "test", "password": "test"}`,
			inputUser: todo.User{
				Name:     "Test",
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user todo.User) {
				r.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Send test request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.SignInInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "test", "password": "test"}`,
			inputUser: todo.SignInInput{
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user todo.SignInInput) {
				r.EXPECT().GenerateToken(user.UserName, user.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Empty fileds",
			inputBody:            `{"username": "test"}`,
			mockBehavior:         func(r *mock_service.MockAuthorization, user todo.SignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"username": "test", "password": "test"}`,
			inputUser: todo.SignInInput{
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user todo.SignInInput) {
				r.EXPECT().GenerateToken(user.UserName, user.Password).Return("token", errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			service := &service.Service{
				Authorization: auth,
			}
			handler := NewHandler(service)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
