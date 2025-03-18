package dogs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/models"
	"github.com/rhargreaves/dog-walking/api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireAPIError(t *testing.T, c *gin.Context, expectedCode int, expectedMessage string) {
	require.Len(t, c.Errors, 1)
	if apiError, ok := c.Errors[0].Err.(common.APIError); ok {
		assert.Equal(t, expectedCode, apiError.Code)
		require.Equal(t, expectedMessage, apiError.Message)
	} else {
		t.Fatalf("Expected APIError, got %T", c.Errors[0].Err)
	}
}

func TestUploadPhoto_ReturnsBadRequest_WhenFileIsNotAnImage(t *testing.T) {
	const dogId = "123"

	dogRepository := new(mocks.DogRepository)
	dogRepository.EXPECT().Get(dogId).Return(&models.Dog{ID: dogId}, nil)
	handler := NewDogPhotoHandler(dogRepository, nil, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: dogId}}

	c.Request = httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/dogs/%s/photo", dogId), strings.NewReader("not an image"))
	c.Request.Header.Set("Content-Type", "text/plain")

	handler.UploadDogPhoto(c)

	requireAPIError(t, c, http.StatusBadRequest, "invalid image content type")
}

func TestUploadPhoto_ReturnsBadRequest_WhenImageIsNotADog(t *testing.T) {
	const dogId = "123"

	dogRepository := new(mocks.DogRepository)
	dogRepository.EXPECT().Get(dogId).Return(&models.Dog{ID: dogId}, nil)
	breedDetector := new(mocks.BreedDetector)
	breedDetector.EXPECT().DetectBreed(dogId).Return("", 0.0, ErrNoDogDetected)
	handler := NewDogPhotoHandler(dogRepository, nil, breedDetector)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: dogId}}

	c.Request = httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/dogs/%s/photo/detect-breed", dogId), strings.NewReader("dummy"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.DetectBreed(c)

	requireAPIError(t, c, http.StatusBadRequest, "no dog detected")
}

func TestUploadPhoto_ReturnsBadRequest_WhenNoSpecificBreedOfDogIsDetected(t *testing.T) {
	const dogId = "123"

	dogRepository := new(mocks.DogRepository)
	dogRepository.EXPECT().Get(dogId).Return(&models.Dog{ID: dogId}, nil)
	breedDetector := new(mocks.BreedDetector)
	breedDetector.EXPECT().DetectBreed(dogId).Return("", 0.0, ErrNoSpecificBreedDetected)
	handler := NewDogPhotoHandler(dogRepository, nil, breedDetector)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: dogId}}

	c.Request = httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/dogs/%s/photo/detect-breed", dogId), strings.NewReader("dummy"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.DetectBreed(c)

	requireAPIError(t, c, http.StatusBadRequest, "no specific breed detected")
}
