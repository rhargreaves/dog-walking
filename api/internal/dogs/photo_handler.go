package dogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DogPhotoHandler interface {
	UploadDogPhoto(c *gin.Context)
}

type dogPhotoHandler struct {
	dogRepository      DogRepository
	dogPhotoRepository DogPhotoRepository
}

func NewDogPhotoHandler(dogRepository DogRepository, dogPhotoRepository DogPhotoRepository) DogPhotoHandler {
	return &dogPhotoHandler{
		dogRepository:      dogRepository,
		dogPhotoRepository: dogPhotoRepository,
	}
}

func (h *dogPhotoHandler) UploadDogPhoto(c *gin.Context) {
	id := c.Param("id")
	_, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	err = h.dogPhotoRepository.Upload(id, c.Request.Body, c.GetHeader("Content-Type"))
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
