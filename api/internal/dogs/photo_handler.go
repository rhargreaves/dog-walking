package dogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogPhotoHandler interface {
	UploadDogPhoto(c *gin.Context)
}

type dogPhotoHandler struct {
	dogRepository    DogRepository
	dogPhotoUploader DogPhotoUploader
}

func NewDogPhotoHandler(dogRepository DogRepository, dogPhotoUploader DogPhotoUploader) DogPhotoHandler {
	return &dogPhotoHandler{
		dogRepository:    dogRepository,
		dogPhotoUploader: dogPhotoUploader,
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

	_, err = h.dogPhotoUploader.Upload(id, c.Request.Body, contentType)
	if err != nil {
		handleError(c, err)
		return
	}

	err = h.dogRepository.UpdatePhotoStatus(id, "pending")
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
