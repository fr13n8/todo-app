package handler

import (
	"net/http"
	"strconv"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type getAllListResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getAllListResponse{
		Data: lists,
	})
}

type getListResponse struct {
	Data todo.TodoList `json:"data"`
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	list, err := h.services.TodoList.GetById(listId, userId)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getListResponse{
		Data: list,
	})
}

func (h *Handler) updateList(c *gin.Context) {
	var input todo.UpdateListInput
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.TodoList.Update(listId, userId, input)

	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.TodoList.Delete(listId, userId)

	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
