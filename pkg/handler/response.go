package handler

import (
	"github.com/gin-gonic/gin"
)

func newReponseError(c *gin.Context, status int, err error) {
	er := HTTPError{
		Message: err.Error(),
	}
	c.AbortWithStatusJSON(status, er)
}

type HTTPError struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}
