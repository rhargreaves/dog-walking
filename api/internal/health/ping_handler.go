package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler godoc
// @Summary Health check endpoint
// @Description Returns OK if the API is running
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
