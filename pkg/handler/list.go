package handler

import (
	"errors"
	"net/http"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		newReponseError(c, http.StatusUnauthorized, errors.New("user id not found"))
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, err)
	}

}

func (h *Handler) getAllList(c *gin.Context) {
}

func (h *Handler) getListById(c *gin.Context) {

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
