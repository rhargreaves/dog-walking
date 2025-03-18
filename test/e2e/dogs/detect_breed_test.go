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

func uploadImageAndDetectBreed(t *testing.T, dogID string, imagePath string) *http.Response {
	image, err := os.ReadFile(imagePath)
	require.NoError(t, err)

	resp := putBytes(t, fmt.Sprintf("/dogs/%s/photo", dogID),
		image, "image/jpeg")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	return common.PostJson(t,
		fmt.Sprintf("/dogs/%s/photo/detect-breed", dogID),
		DetectBreedRequest{})
}

func TestDetectBreed_SuccessfulCases(t *testing.T) {
	tests := []breedTestCase{
		{
			dogName:       "Mr.Peanutbutter",
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
			dog := createDog(t, tc.dogName)
			resp := uploadImageAndDetectBreed(t, dog.ID, tc.imagePath)
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusOK)

			var response DetectBreedResponse
			common.DecodeJSON(t, resp, &response)
			require.Equal(t, tc.expectedBreed, response.Breed)
			require.Greater(t, response.Confidence, 55.0)

			resp = common.Get(t, "/dogs/"+dog.ID)
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusOK)

			var fetchedDog Dog
			common.DecodeJSON(t, resp, &fetchedDog)
			require.Equal(t, tc.expectedBreed, fetchedDog.Breed)
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
			dog := createDog(t, tc.dogName)
			resp := uploadImageAndDetectBreed(t, dog.ID, tc.imagePath)
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusBadRequest)

			var response common.ErrorResponse
			common.DecodeJSON(t, resp, &response)
			require.Equal(t, "no dog detected", response.Error)
		})
	}
}
