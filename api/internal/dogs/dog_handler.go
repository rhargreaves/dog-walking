package dogs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/models"
)

type DogHandler interface {
	ListDogs(c *gin.Context)
	GetDog(c *gin.Context)
	CreateDog(c *gin.Context)
	UpdateDog(c *gin.Context)
	DeleteDog(c *gin.Context)
}

type dogHandler struct {
	dogRepository DogRepository
}

func NewDogHandler(dogRepository DogRepository) DogHandler {
	return &dogHandler{dogRepository: dogRepository}
}

func (h *dogHandler) CreateDog(c *gin.Context) {
	var dog models.Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Create(&dog); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dog)
}

func (h *dogHandler) ListDogs(c *gin.Context) {
	dogs, err := h.dogRepository.List()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dogs)
}

func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) UpdateDog(c *gin.Context) {
	id := c.Param("id")
	var dog models.Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Update(id, &dog); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) DeleteDog(c *gin.Context) {
	id := c.Param("id")
	if err := h.dogRepository.Delete(id); err != nil {
		handleError(c, err)
		return
	}
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, ErrDogNotFound) {
		c.Error(common.APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	c.Error(err)
}

func handleBindError(c *gin.Context, err error) {
	c.Error(common.APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}
