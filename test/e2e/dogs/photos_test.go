package dogs

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
	"github.com/stretchr/testify/require"
)

func TestUploadImage_PhotoUploadedToS3(t *testing.T) {
	dog := createDog(t, "Mr. Peanutbutter")

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)

	req := common.ApiRequest(t, http.MethodPut, fmt.Sprintf("/dogs/%s/photo", dog.ID),
		true, bytes.NewReader(image))
	req.Header.Set("Content-Type", "image/jpeg")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
	t.Log("Image uploaded successfully")

	s3Image := getS3Object(t, os.Getenv("DOG_IMAGES_BUCKET"), dog.ID)

	require.Equal(t, image, s3Image)
}

func TestUploadImage_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	req := common.ApiRequest(t, http.MethodPut, "/dogs/123/photo",
		true, bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "text/plain")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}
