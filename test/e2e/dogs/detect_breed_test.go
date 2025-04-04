package dogs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func uploadImageAndGetDog(t *testing.T, dogID string, imagePath string) *DogResponse {
	image, err := os.ReadFile(imagePath)
	require.NoError(t, err)

	uploadImage(t, dogID, image)

	return getDogWaitingForPhotoModeration(t, dogID)
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
			dog := createDog(t, CreateOrUpdateDogRequest{Name: tc.dogName})
			fetchedDog := uploadImageAndGetDog(t, dog.ID, tc.imagePath)

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
			dog := createDog(t, CreateOrUpdateDogRequest{Name: tc.dogName})
			fetchedDog := uploadImageAndGetDog(t, dog.ID, tc.imagePath)
			require.Equal(t, "rejected", fetchedDog.PhotoStatus)
		})
	}
}
