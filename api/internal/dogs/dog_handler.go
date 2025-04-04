package dogs

import (
	"errors"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/model"
)

type DogHandler interface {
	ListDogs(c *gin.Context)
	GetDog(c *gin.Context)
	CreateDog(c *gin.Context)
	UpdateDog(c *gin.Context)
	DeleteDog(c *gin.Context)
}

type dogHandler struct {
	config        *DogHandlerConfig
	dogRepository domain.DogRepository
}

type DogHandlerConfig struct {
	ImagesCdnBaseUrl string
}

func NewDogHandler(dogHandlerConfig DogHandlerConfig, dogRepository domain.DogRepository) DogHandler {
	return &dogHandler{config: &dogHandlerConfig, dogRepository: dogRepository}
}

// CreateDog godoc
// @Summary Create a new dog
// @Description Create a new dog with the provided details
// @Tags dogs
// @Accept json
// @Produce json
// @Param dog body model.CreateOrUpdateDogRequest true "Dog information"
// @Success 201 {object} model.DogResponse
// @Failure 400 {object} common.APIErrorResponse "Invalid request"
// @Failure 500 {object} common.APIErrorResponse "Internal server error"
// @Router /dogs [post]
func (h *dogHandler) CreateDog(c *gin.Context) {
	var request model.CreateOrUpdateDogRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		handleBindError(c, err)
		return
	}

	dog := *model.FromCreateOrUpdateDogRequest(&request)
	createdDog, err := h.dogRepository.Create(dog)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.ToDogResponse(createdDog, h.config.ImagesCdnBaseUrl))
}

// ListDogs godoc
// @Summary List all dogs
// @Description Get a list of all registered dogs
// @Tags dogs
// @Produce json
// @Param limit query int false "Limit the number of dogs returned" default(25) minimum(1) maximum(25)
// @Param name query string false "Filter dogs by partial name"
// @Param nextToken query string false "A token to get the next page of results"
// @Success 200 {object} model.DogListResponse
// @Failure 500 {object} common.APIErrorResponse "Internal server error"
// @Router /dogs [get]
func (h *dogHandler) ListDogs(c *gin.Context) {
	var query model.DogListQuery
	defaults.Set(&query)

	if err := c.ShouldBindQuery(&query); err != nil {
		handleBindError(c, err)
		return
	}

	dogs, err := h.dogRepository.List(query.Limit, query.Name, query.NextToken)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToDogListResponse(dogs, h.config.ImagesCdnBaseUrl))
}

// GetDog godoc
// @Summary Get a dog by ID
// @Description Get details of a specific dog by its ID
// @Tags dogs
// @Produce json
// @Param id path string true "Dog ID"
// @Success 200 {object} model.DogResponse
// @Failure 404 {object} common.APIErrorResponse "Dog not found"
// @Failure 500 {object} common.APIErrorResponse "Internal server error"
// @Router /dogs/{id} [get]
func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToDogResponse(dog, h.config.ImagesCdnBaseUrl))
}

// UpdateDog godoc
// @Summary Update a dog
// @Description Update a dog's information by its ID
// @Tags dogs
// @Accept json
// @Produce json
// @Param id path string true "Dog ID"
// @Param dog body model.CreateOrUpdateDogRequest true "Updated dog information"
// @Success 200 {object} model.DogResponse
// @Failure 400 {object} common.APIErrorResponse "Invalid request"
// @Failure 404 {object} common.APIErrorResponse "Dog not found"
// @Failure 500 {object} common.APIErrorResponse "Internal server error"
// @Router /dogs/{id} [put]
func (h *dogHandler) UpdateDog(c *gin.Context) {
	id := c.Param("id")
	var request model.CreateOrUpdateDogRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		handleBindError(c, err)
		return
	}

	dog := *model.FromCreateOrUpdateDogRequest(&request)
	if err := h.dogRepository.Update(id, &dog); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToDogResponse(&dog, h.config.ImagesCdnBaseUrl))
}

// DeleteDog godoc
// @Summary Delete a dog
// @Description Delete a dog by its ID
// @Tags dogs
// @Param id path string true "Dog ID"
// @Success 204 "No Content"
// @Failure 404 {object} common.APIErrorResponse "Dog not found"
// @Failure 500 {object} common.APIErrorResponse "Internal server error"
// @Router /dogs/{id} [delete]
func (h *dogHandler) DeleteDog(c *gin.Context) {
	id := c.Param("id")
	if err := h.dogRepository.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrDogNotFound) {
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
