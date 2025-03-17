package dogs

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
	"github.com/stretchr/testify/require"
)

const testCartoonDogImagePath = "../resources/mr_peanutbutter.jpg"
const testToyImagePath = "../resources/toy.jpg"
const testCatImagePath = "../resources/cat.jpg"
const testHuskyImagePath = "../resources/husky.jpg"

type DetectBreedRequest struct {
}

type DetectBreedResponse struct {
	Breed      string  `json:"breed"`
	Confidence float64 `json:"confidence"`
}

func TestUploadImage_PhotoUploadedToS3(t *testing.T) {
	dog := createDog(t, "Mr. Peanutbutter")

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
	t.Log("Image uploaded successfully")

	s3Image := getS3Object(t, os.Getenv("DOG_IMAGES_BUCKET"), dog.ID)

	require.Equal(t, image, s3Image)
}

func TestUploadImage_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), "123"),
		[]byte{}, "text/plain")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}

func TestDetectBreed_PopulatesBreedAttribute(t *testing.T) {
	dog := createDog(t, "Mr. Peanutbutter")

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.PostJson(t, fmt.Sprintf("/dogs/%s/photo/detect-breed", dog.ID),
		DetectBreedRequest{})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var response DetectBreedResponse
	common.DecodeJSON(t, resp, &response)

	require.Equal(t, "Airedale", response.Breed)
	require.Greater(t, response.Confidence, 55.0)
}

func TestDetectBreed_ReturnsHuskyBreed(t *testing.T) {
	dog := createDog(t, "Husky")

	image, err := os.ReadFile(testHuskyImagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.PostJson(t, fmt.Sprintf("/dogs/%s/photo/detect-breed", dog.ID),
		DetectBreedRequest{})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var response DetectBreedResponse
	common.DecodeJSON(t, resp, &response)

	require.Equal(t, "Husky", response.Breed)
	require.Greater(t, response.Confidence, 55.0)
}

func TestDetectBreed_ReturnsNotADogForNotAnAnimal(t *testing.T) {
	dog := createDog(t, "Sweep")

	image, err := os.ReadFile(testToyImagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.PostJson(t, fmt.Sprintf("/dogs/%s/photo/detect-breed", dog.ID),
		DetectBreedRequest{})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusBadRequest)

	var response common.ErrorResponse
	common.DecodeJSON(t, resp, &response)

	require.Equal(t, "no dog detected", response.Error)
}

func TestDetectBreed_ReturnsNotADogForCat(t *testing.T) {
	dog := createDog(t, "Cat")

	image, err := os.ReadFile(testCatImagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.PostJson(t, fmt.Sprintf("/dogs/%s/photo/detect-breed", dog.ID),
		DetectBreedRequest{})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusBadRequest)

	var response common.ErrorResponse
	common.DecodeJSON(t, resp, &response)

	require.Equal(t, "no dog detected", response.Error)
}
