package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/fr13n8/todo-app/pkg/service"
	mockservice "github.com/fr13n8/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(&jwt.StandardClaims{
					Issuer:    "",
					ExpiresAt: 0,
					IssuedAt:  0,
					Id:        "1",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `1`,
		},
		{
			name:                 "No Header",
			headerName:           "Authorization",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth token header"}`,
		},
		{
			name:                 "Invalid header",
			headerName:           "Authorization",
			headerValue:          "Berer token",
			token:                "token",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"broken auth token"}`,
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"broken auth token"}`,
		},
		{
			name:        "Service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(nil, errors.New("failed to parse token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
