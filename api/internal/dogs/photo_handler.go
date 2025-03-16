package dogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DogPhotoHandler interface {
	UploadDogPhoto(c *gin.Context)
	DetectBreed(c *gin.Context)
}

type dogPhotoHandler struct {
	dogRepository      DogRepository
	dogPhotoRepository DogPhotoRepository
	breedDetector      BreedDetector
}

func NewDogPhotoHandler(dogRepository DogRepository, dogPhotoRepository DogPhotoRepository, breedDetector BreedDetector) DogPhotoHandler {
	return &dogPhotoHandler{
		dogRepository:      dogRepository,
		dogPhotoRepository: dogPhotoRepository,
		breedDetector:      breedDetector,
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

func (h *dogPhotoHandler) DetectBreed(c *gin.Context) {
	id := c.Param("id")

	// Check if dog exists and get current dog data
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	breed, confidence, err := h.breedDetector.DetectBreed(id)
	if err != nil {
		handleError(c, err)
		return
	}

	// Update the dog's breed in the database
	dog.Breed = breed
	err = h.dogRepository.Update(id, dog)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"breed":      breed,
		"confidence": confidence,
	})
}
