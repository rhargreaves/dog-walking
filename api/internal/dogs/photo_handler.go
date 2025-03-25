package dogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/model"
)

type DogPhotoHandler interface {
	UploadDogPhoto(c *gin.Context)
	DetectBreed(c *gin.Context)
}

type dogPhotoHandler struct {
	dogRepository    DogRepository
	dogPhotoUploader DogPhotoUploader
	breedDetector    BreedDetector
}

func NewDogPhotoHandler(dogRepository DogRepository, dogPhotoUploader DogPhotoUploader, breedDetector BreedDetector) DogPhotoHandler {
	return &dogPhotoHandler{
		dogRepository:    dogRepository,
		dogPhotoUploader: dogPhotoUploader,
		breedDetector:    breedDetector,
	}
}

// UploadDogPhoto godoc
// @Summary Upload a dog's photo
// @Description Upload a JPEG photo for a specific dog
// @Tags dogs,photos
// @Accept image/jpeg
// @Param id path string true "Dog ID"
// @Success 200 "OK"
// @Failure 400 {object} common.APIError "Invalid content type or request"
// @Failure 404 {object} common.APIError "Dog not found"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs/{id}/photo [put]
func (h *dogPhotoHandler) UploadDogPhoto(c *gin.Context) {
	id := c.Param("id")
	_, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	contentType := c.GetHeader("Content-Type")
	if contentType != "image/jpeg" {
		c.Error(common.APIError{
			Code:    http.StatusBadRequest,
			Message: "invalid image content type",
		})
		return
	}

	err = h.dogPhotoUploader.Upload(id, c.Request.Body, contentType)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// DetectBreed godoc
// @Summary Detect a dog's breed from its photo
// @Description Analyzes a previously uploaded photo to detect the dog's breed
// @Tags dogs,photos
// @Param id path string true "Dog ID"
// @Produce json
// @Success 200 {object} model.BreedDetectionResultResponse "Returns id, breed, and confidence"
// @Failure 400 {object} common.APIError "No dog detected or no specific breed detected"
// @Failure 404 {object} common.APIError "Dog not found"
// @Failure 500 {object} common.APIError "Internal server error"
// @Router /dogs/{id}/photo/detect-breed [post]
func (h *dogPhotoHandler) DetectBreed(c *gin.Context) {
	id := c.Param("id")

	// Check if dog exists and get current dog data
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	breedResult, err := h.breedDetector.DetectBreed(id)
	if err != nil {
		if err == ErrNoDogDetected || err == ErrNoSpecificBreedDetected {
			c.Error(common.APIError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		handleError(c, err)
		return
	}

	// Update the dog's breed in the database
	dog.Breed = breedResult.Breed
	err = h.dogRepository.Update(id, dog)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToBreedDetectionResultResponse(id, breedResult))
}
