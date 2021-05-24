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
		newReponseError(c, http.StatusUnauthorized, errors.New("empty auth token header"))
		return
	}

	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		newReponseError(c, http.StatusUnauthorized, errors.New("Broken auth token!"))
		return
	}
	token := splitToken[1]

	userId, err := h.services.ParseToken(token)
	if err != nil {
		newReponseError(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtx, userId)

}
