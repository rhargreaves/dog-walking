package dogs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireAPIError(t *testing.T, c *gin.Context, expectedCode int, expectedMessage string) {
	require.Len(t, c.Errors, 1)
	if apiError, ok := c.Errors[0].Err.(common.APIError); ok {
		assert.Equal(t, expectedCode, apiError.Code)
		require.Equal(t, expectedMessage, apiError.Message)
	} else {
		t.Fatalf("Expected APIError, got %T", c.Errors[0].Err)
	}
}

func TestUploadPhoto_ReturnsBadRequest_WhenFileIsNotAnImage(t *testing.T) {
	dogRepository := NewFakeDogRepository()
	dog, _ := dogRepository.Create(domain.Dog{})
	handler := NewDogPhotoHandler(dogRepository, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: dog.ID}}

	c.Request = httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/dogs/%s/photo", dog.ID), strings.NewReader("not an image"))
	c.Request.Header.Set("Content-Type", "text/plain")

	handler.UploadDogPhoto(c)

	requireAPIError(t, c, http.StatusBadRequest, "invalid image content type")
}
