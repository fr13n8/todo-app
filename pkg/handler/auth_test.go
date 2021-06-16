package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/service"
	mockservice "github.com/fr13n8/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, user todo.SignUpInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.SignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "username": "test", "password": "test"}`,
			inputUser: todo.SignUpInput{
				Name:     "Test",
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mockservice.MockAuthorization, user todo.SignUpInput) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"username": "test", "password": "test"}`,
			mockBehavior:         func(r *mockservice.MockAuthorization, user todo.SignUpInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name": "Test", "username": "test", "password": "test"}`,
			inputUser: todo.SignUpInput{
				Name:     "Test",
				UserName: "test",
				Password: "test",
			},
			mockBehavior: func(r *mockservice.MockAuthorization, user todo.SignUpInput) {
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

			auth := mockservice.NewMockAuthorization(c)
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
	type input struct {
		user      todo.SignInInput
		userAgent string
	}

	type mockBehavior func(s *mockservice.MockAuthorization, input input)

	testTable := []struct {
		name                 string
		inputBody            string
		input                input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "test", "password": "test"}`,
			input: input{
				user: todo.SignInInput{
					UserName: "test",
					Password: "test",
				},
				userAgent: "test",
			},
			mockBehavior: func(r *mockservice.MockAuthorization, input input) {
				r.EXPECT().SignInUser(input.user.UserName, input.user.Password, input.userAgent).Return([]string{"token", "token1"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"token","refreshToken":"token1"}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"username": "test"}`,
			mockBehavior:         func(r *mockservice.MockAuthorization, input input) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "service failure",
			inputBody: `{"username": "test", "password": "test"}`,
			input: input{
				user: todo.SignInInput{
					UserName: "test",
					Password: "test",
				},
				userAgent: "test",
			},
			mockBehavior: func(r *mockservice.MockAuthorization, input input) {
				r.EXPECT().SignInUser(input.user.UserName, input.user.Password, input.userAgent).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.input)

			s := &service.Service{Authorization: auth}
			handler := NewHandler(s)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))
			req.Header.Add("User-Agent", testCase.input.userAgent)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_refreshToken(t *testing.T) {
	type input todo.RefreshTokenInput

	type mockBehavior func(s *mockservice.MockAuthorization, input input)

	testTable := []struct {
		name                 string
		inputBody            string
		input                input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"refresh_token":"token"}`,
			input: input{
				RefreshToken: "token",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, input input) {
				s.EXPECT().RefreshToken(input.RefreshToken).Return([]string{"token", "token1"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"token","refreshToken":"token1"}`,
		},
		{
			name:      "Expired token",
			inputBody: `{"refresh_token":"token"}`,
			input: input{
				RefreshToken: "token",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, input input) {
				s.EXPECT().RefreshToken(input.RefreshToken).Return(nil, errors.New("refresh token expired"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"refresh token expired"}`,
		},
		{
			name:                 "Without token",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockservice.MockAuthorization, input input) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.input)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			r := gin.New()
			r.POST("/refresh", handler.refreshToken)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
