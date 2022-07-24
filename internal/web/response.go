package web

import (
	"github.com/gin-gonic/gin"
)

type Respons struct {
	Status  int         `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponError struct {
	Status  int         `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func MarshalPayload(c *gin.Context, code int, message string, payload interface{}) {
	res := Respons{
		Status:  code,
		Message: message,
		Data:    payload,
	}

	c.JSON(code, res)
}

func MarshalError(c *gin.Context, code int, message string, payload interface{}) {
	res := ResponError{
		Status:  code,
		Message: message,
		Data:    payload,
	}

	c.JSON(code, res)
}
