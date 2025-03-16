package dogs

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
	"github.com/stretchr/testify/require"
)

func TestUploadImage(t *testing.T) {
	dog := createDog(t, "Rover")

	file, err := os.Open("../resources/dog.jpg")
	require.NoError(t, err)
	defer file.Close()

	url := fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID)
	req, err := http.NewRequest("PUT", url, file)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	t.Log("Image uploaded successfully")

	requireS3ObjectExists(t, os.Getenv("DOG_IMAGES_BUCKET"), dog.ID)
}
