package dogs

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func dogWithPhotoUrl(dog *models.Dog) *models.Dog {
	if dog.PhotoHash != "" {
		dog.PhotoUrl = fmt.Sprintf("%s/%s?h=%s",
			os.Getenv("CLOUDFRONT_BASE_URL"), dog.ID, dog.PhotoHash)
	}
	return dog
}

type dogHandler struct {
	dogRepository DogRepository
}

func NewDogHandler(dogRepository DogRepository) DogHandler {
	return &dogHandler{dogRepository: dogRepository}
}

// CreateDog godoc
// @Summary Create a new dog
// @Description Create a new dog with the provided details
// @Tags dogs
// @Accept json
// @Produce json
// @Param dog body models.Dog true "Dog information"
// @Success 201 {object} models.Dog
// @Failure 400 {object} common.APIError "Invalid request"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs [post]
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

// ListDogs godoc
// @Summary List all dogs
// @Description Get a list of all registered dogs
// @Tags dogs
// @Produce json
// @Success 200 {array} models.Dog
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs [get]
func (h *dogHandler) ListDogs(c *gin.Context) {
	limit := c.Query("limit")
	nextToken := c.Query("nextToken")

	var limitInt int = 25
	if limit != "" {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			handleError(c, err)
			return
		}
	}

	dogs, err := h.dogRepository.List(limitInt, nextToken)
	if err != nil {
		handleError(c, err)
		return
	}

	for i := range dogs.Dogs {
		dogs.Dogs[i] = *dogWithPhotoUrl(&dogs.Dogs[i])
	}

	c.JSON(http.StatusOK, dogs)
}

// GetDog godoc
// @Summary Get a dog by ID
// @Description Get details of a specific dog by its ID
// @Tags dogs
// @Produce json
// @Param id path string true "Dog ID"
// @Success 200 {object} models.Dog
// @Failure 404 {object} common.APIError "Dog not found"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs/{id} [get]
func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dogWithPhotoUrl(dog))
}

// UpdateDog godoc
// @Summary Update a dog
// @Description Update a dog's information by its ID
// @Tags dogs
// @Accept json
// @Produce json
// @Param id path string true "Dog ID"
// @Param dog body models.Dog true "Updated dog information"
// @Success 200 {object} models.Dog
// @Failure 400 {object} common.APIError "Invalid request"
// @Failure 404 {object} common.APIError "Dog not found"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs/{id} [put]
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

// DeleteDog godoc
// @Summary Delete a dog
// @Description Delete a dog by its ID
// @Tags dogs
// @Param id path string true "Dog ID"
// @Success 204 "No Content"
// @Failure 404 {object} common.APIError "Dog not found"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs/{id} [delete]
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
