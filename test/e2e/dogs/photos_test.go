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
	dog := createDog(t, "Mr. Peanutbutter")

	image, err := os.ReadFile(testCartoonDogImagePath)
	require.NoError(t, err)
	uploadImage(t, dog.ID, image)

	cdnUrl := fmt.Sprintf("%s/%s", os.Getenv("CLOUDFRONT_BASE_URL"), dog.ID)
	req, err := http.NewRequest(http.MethodGet, cdnUrl, nil)
	require.NoError(t, err)
	resp := common.GetResponse(t, req)
	assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
	t.Log("Image accessible via CDN")
}
