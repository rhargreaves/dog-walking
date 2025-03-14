package dogs

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

func PostDog(c *gin.Context) {
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		c.Error(common.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	r := NewDogRepository(os.Getenv("DOGS_TABLE_NAME"))
	err := r.Create(&dog)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func ListDogs(c *gin.Context) {
	r := NewDogRepository(os.Getenv("DOGS_TABLE_NAME"))
	dogs, err := r.List()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dogs)
}
