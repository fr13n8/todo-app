package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

// @Summary Create todo item
// @Security ApiKeyAuth
// @Tags items
// @Description create todo item
// @ID create-item
// @Accept  json
// @Produce  json
// @Param input body todo.Item true "item info"
// @Param id path int true "list id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	var input todo.Item
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.services.TodoItem.Create(listId, userId, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type getAllItemsResponse struct {
	Data []todo.Item `json:"data"`
}

// @Summary Get All items
// @Security ApiKeyAuth
// @Tags items
// @Description get all items
// @ID get-all-items
// @Accept  json
// @Produce  json
// @Param id path int true "list id"
// @Success 200 {object} getAllItemsResponse
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	items, err := h.services.TodoItem.GetAll(listId, userId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

type getItemResponse struct {
	Data todo.Item `json:"data"`
}

// @Summary Get item by id
// @Security ApiKeyAuth
// @Tags items
// @Description get item by id
// @ID get-item-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "list id"
// @Param item_id path int true "item id"
// @Success 200 {object} getItemResponse
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id/items/:item_id [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getItemResponse{
		Data: item,
	})
}

// @Summary Update todo item
// @Security ApiKeyAuth
// @Tags items
// @Description update todo item
// @ID update-item
// @Accept  json
// @Produce  json
// @Param id path int true "item id"
// @Param input body todo.UpdateItemInput true "item info"
// @Success 200 {string} string Ok
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/items/:id [put]
func (h *Handler) updateItem(c *gin.Context) {
	var input todo.UpdateItemInput
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := h.services.TodoItem.Update(userId, itemId, input); err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

// @Summary Delete item
// @Security ApiKeyAuth
// @Tags items
// @Description delete item by id
// @ID delete-item-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "list id"
// @Param item_id path int true "item id"
// @Success 200 {string} Ok
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id/items/:item_id [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
