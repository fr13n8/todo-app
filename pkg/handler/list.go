package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fr13n8/todo-app"
	"github.com/gin-gonic/gin"
)

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body todo.List true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input todo.List
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type getAllListResponse struct {
	Data []todo.List `json:"data"`
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListResponse
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists [get]
func (h *Handler) getAllList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getAllListResponse{
		Data: lists,
	})
}

type getListResponse struct {
	Data todo.List `json:"data"`
}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "List id"
// @Success 200 {object} getListResponse
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	list, err := h.services.TodoList.GetById(listId, userId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getListResponse{
		Data: list,
	})
}

// @Summary Update todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description update todo list
// @ID update-list
// @Accept  json
// @Produce  json
// @Param input body todo.UpdateListInput true "list info"
// @Success 200 {string} string Ok
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists [put]
func (h *Handler) updateList(c *gin.Context) {
	var input todo.UpdateListInput
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := h.services.TodoList.Update(listId, userId, input); err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

// @Summary Delete List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description delete list by id
// @ID delete-list-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "List id"
// @Success 200 {string} Ok
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /api/lists/:id [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.TodoList.Delete(listId, userId)

	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
