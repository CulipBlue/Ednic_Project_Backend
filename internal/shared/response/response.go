package response

import "github.com/gin-gonic/gin"

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(200, Body{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

func Error(c *gin.Context, status int, message string, errors any) {
	c.JSON(status, Body{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
	})
}
