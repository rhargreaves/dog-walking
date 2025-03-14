package dogs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogHandler interface {
	CreateDog(c *gin.Context)
	ListDogs(c *gin.Context)
	GetDog(c *gin.Context)
	UpdateDog(c *gin.Context)
}

type dogHandler struct {
	dogRepository DogRepository
}

func NewDogHandler(dogRepository DogRepository) DogHandler {
	return &dogHandler{dogRepository: dogRepository}
}

func (h *dogHandler) CreateDog(c *gin.Context) {
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		c.Error(common.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := h.dogRepository.Create(&dog)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) ListDogs(c *gin.Context) {
	dogs, err := h.dogRepository.List()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dogs)
}

func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		if errors.Is(err, ErrDogNotFound) {
			c.Error(common.APIError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		} else {
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) UpdateDog(c *gin.Context) {
	id := c.Param("id")
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		c.Error(common.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := h.dogRepository.Update(id, &dog)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dog)
}
