package common

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

type APIErrorResponse struct {
	Error APIError `json:"error"`
}

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		for _, e := range c.Errors {
			if apiErr, ok := e.Err.(APIError); ok {
				c.JSON(apiErr.Code, APIErrorResponse{Error: apiErr})
				return
			}
			c.JSON(http.StatusInternalServerError, APIErrorResponse{Error: APIError{
				Code:    http.StatusInternalServerError,
				Message: e.Error(),
			}})
			return
		}
	}
}
