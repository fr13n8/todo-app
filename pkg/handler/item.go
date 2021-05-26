package handler

import (
	"net/http"
	"strconv"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.TodoItem.Create(listId, userId, input)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type getAllItemsResponse struct {
	Data []todo.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	items, err := h.services.TodoItem.GetAll(listId, userId)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

type getItemResponse struct {
	Data todo.TodoItem `json:"data"`
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getItemResponse{
		Data: item,
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	var input todo.UpdateItemInput
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.services.TodoItem.Update(userId, itemId, input); err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newReponseError(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newReponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
