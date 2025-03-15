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
		h.handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Create(&dog); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) ListDogs(c *gin.Context) {
	dogs, err := h.dogRepository.List()
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dogs)
}

func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) UpdateDog(c *gin.Context) {
	id := c.Param("id")
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		h.handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Update(id, &dog); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, ErrDogNotFound) {
		c.Error(common.APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	c.Error(err)
}

func (h *dogHandler) handleBindError(c *gin.Context, err error) {
	c.Error(common.APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}
