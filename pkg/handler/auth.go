package handler

import (
	"errors"
	"net/http"

	"github.com/fr13n8/todo-app/structs"
	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body structs.SignUpInput true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input structs.SignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary SignIn
// @Tags auth
// @Description user login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body structs.SignInInput true "credentials"
// @Success 200 {object} AuthResponse
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input structs.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}
	userAgent := c.Request.Header.Get("User-Agent")
	tokens, err := h.services.SignInUser(input.UserName, input.Password, userAgent)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	})
}

// @Summary Refresh
// @Tags auth
// @Description refresh JWT token
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body structs.RefreshTokenInput true "refresh token"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /auth/refresh [post]
func (h *Handler) refreshToken(c *gin.Context) {
	var input structs.RefreshTokenInput

	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	tokens, err := h.services.RefreshToken(input.RefreshToken)
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	})
}
