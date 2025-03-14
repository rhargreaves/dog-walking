package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func errorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		for _, e := range c.Errors {
			if apiErr, ok := e.Err.(APIError); ok {
				c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			return
		}
	}
}
