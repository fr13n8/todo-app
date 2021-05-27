package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newResponseError(c, http.StatusUnauthorized, errors.New("empty auth token header"))
		return
	}

	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		newResponseError(c, http.StatusUnauthorized, errors.New("broken auth token"))
		return
	}
	token := splitToken[1]

	userId, err := h.services.ParseToken(token)
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newResponseError(c, http.StatusInternalServerError, errors.New("user id not found"))
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newResponseError(c, http.StatusInternalServerError, errors.New("invalid user id type"))
		return 0, errors.New("invalid user id type")
	}

	return idInt, nil
}
