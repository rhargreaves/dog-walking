package dogs

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uploadImage(t *testing.T, dogID string, image []byte) {
	req := common.ApiRequest(t, http.MethodPut, fmt.Sprintf("/dogs/%s/photo", dogID),
		true, bytes.NewReader(image))
	req.Header.Set("Content-Type", "image/jpeg")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
}

func TestUploadImage_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	req := common.ApiRequest(t, http.MethodPut, "/dogs/123/photo",
		true, bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "text/plain")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}

func TestUploadImage_AccessibleViaCDN(t *testing.T) {
	dog := createDog(t, CreateOrUpdateDogRequest{Name: "Mr. Peanutbutter"})

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)
	uploadImage(t, dog.ID, image)

	getDogWaitingForPhotoModeration(t, dog.ID)

	cdnUrl := fmt.Sprintf("%s/%s", os.Getenv("CLOUDFRONT_BASE_URL"), dog.ID)
	req, err := http.NewRequest(http.MethodGet, cdnUrl, nil)
	require.NoError(t, err)
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
}

func TestUploadImage_ImageStatusIsPendingInDogResponse(t *testing.T) {
	dog := createDog(t, CreateOrUpdateDogRequest{Name: "Mr. Peanutbutter"})

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)
	uploadImage(t, dog.ID, image)

	fetchedDog := getDog(t, dog.ID)
	assert.Equal(t, "pending", fetchedDog.PhotoStatus, "Expected photo status to be pending")
	assert.Empty(t, fetchedDog.PhotoUrl, "Expected photo URL to be empty")
	assert.Empty(t, fetchedDog.PhotoHash, "Expected photo hash to be empty")
}

func TestUploadImage_ImageStatusIsApprovedInDogResponse(t *testing.T) {
	dog := createDog(t, CreateOrUpdateDogRequest{Name: "Mr. Peanutbutter"})

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)
	uploadImage(t, dog.ID, image)
	fetchedDog := getDogWaitingForPhotoModeration(t, dog.ID)

	assert.Equal(t, "approved", fetchedDog.PhotoStatus, "Expected photo status to be approved")
	assert.NotEmpty(t, fetchedDog.PhotoUrl, "Expected photo URL to be returned")
	assert.NotEmpty(t, fetchedDog.PhotoHash, "Expected photo hash to be returned")
}

func TestUploadImage_ImageUrlInDogResponse(t *testing.T) {
	dog := createDog(t, CreateOrUpdateDogRequest{Name: "Mr. Peanutbutter"})

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)
	uploadImage(t, dog.ID, image)

	fetchedDog := getDogWaitingForPhotoModeration(t, dog.ID)
	expectedHash := "443d6817146340599232418cfe7ef31b"
	assert.Equal(t, os.Getenv("CLOUDFRONT_BASE_URL")+"/"+dog.ID+"?h="+expectedHash,
		fetchedDog.PhotoUrl, "Expected photo URL to be returned")
	assert.Equal(t, expectedHash, fetchedDog.PhotoHash, "Expected photo hash to be returned")
}
