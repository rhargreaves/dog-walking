package dogs

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
	"github.com/stretchr/testify/require"
)

type DetectBreedRequest struct {
}

type DetectBreedResponse struct {
	Breed      string  `json:"breed"`
	Confidence float64 `json:"confidence"`
}

type breedTestCase struct {
	dogName       string
	imagePath     string
	expectedBreed string
}

type errorTestCase struct {
	dogName   string
	imagePath string
}

const testCartoonDogImagePath = "../resources/mr_peanutbutter.jpg"
const testToyImagePath = "../resources/toy.jpg"
const testCatImagePath = "../resources/cat.jpg"
const testHuskyImagePath = "../resources/husky.jpg"

func uploadImageAndDetectBreed(t *testing.T, dogName string, imagePath string) (*http.Response, string) {
	dog := createDog(t, dogName)

	image, err := os.ReadFile(imagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	return common.PostJson(t, fmt.Sprintf("/dogs/%s/photo/detect-breed", dog.ID),
		DetectBreedRequest{}), dog.ID
}

func TestDetectBreed_SuccessfulCases(t *testing.T) {
	tests := []breedTestCase{
		{
			dogName:       "Mr. Peanutbutter",
			imagePath:     testCartoonDogImagePath,
			expectedBreed: "Airedale",
		},
		{
			dogName:       "Husky",
			imagePath:     testHuskyImagePath,
			expectedBreed: "Husky",
		},
	}

	for _, tc := range tests {
		t.Run(tc.dogName, func(t *testing.T) {
			resp, _ := uploadImageAndDetectBreed(t, tc.dogName, tc.imagePath)
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusOK)

			var response DetectBreedResponse
			common.DecodeJSON(t, resp, &response)
			require.Equal(t, tc.expectedBreed, response.Breed)
			require.Greater(t, response.Confidence, 55.0)
		})
	}
}

func TestDetectBreed_ErrorCases(t *testing.T) {
	tests := []errorTestCase{
		{
			dogName:   "Sweep",
			imagePath: testToyImagePath,
		},
		{
			dogName:   "Cat",
			imagePath: testCatImagePath,
		},
	}

	for _, tc := range tests {
		t.Run(tc.dogName, func(t *testing.T) {
			resp, _ := uploadImageAndDetectBreed(t, tc.dogName, tc.imagePath)
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusBadRequest)

			var response common.ErrorResponse
			common.DecodeJSON(t, resp, &response)
			require.Equal(t, "no dog detected", response.Error)
		})
	}
}
