package dogs

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
	"github.com/stretchr/testify/require"
)

const testImagePath = "../resources/dog.jpg"

func TestUploadImage_PhotoUploadedToS3(t *testing.T) {
	dog := createDog(t, "Rover")

	image, err := os.ReadFile(testImagePath)
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
