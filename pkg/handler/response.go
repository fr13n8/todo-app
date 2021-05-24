package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func newReponseError(c *gin.Context, status int, err error) {
	logrus.Error(err)
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	c.AbortWithStatusJSON(status, er)
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
