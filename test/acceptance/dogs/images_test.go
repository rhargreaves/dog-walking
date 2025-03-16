package dogs

import (
	"fmt"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
	"github.com/stretchr/testify/require"
)

const testImagePath = "../resources/dog.jpg"

func TestUploadImage_PhotoUploadedToS3(t *testing.T) {
	dog := createDog(t, "Rover")

	image, err := os.Open(testImagePath)
	require.NoError(t, err)
	defer image.Close()
	uploadFile(t, fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID),
		image, "image/jpeg")
	t.Log("Image uploaded successfully")

	uploadedImage := getS3Object(t, os.Getenv("DOG_IMAGES_BUCKET"), dog.ID)
	originalImage, err := os.ReadFile(testImagePath)
	require.NoError(t, err)

	require.Equal(t, originalImage, uploadedImage)
}
