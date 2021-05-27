package handler

import (
	"errors"
	"net/http"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ID": id})
}

func (h *Handler) signIn(c *gin.Context) {

	var input todo.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	token, err := h.services.GenerateToken(input.UserName, input.Password)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
