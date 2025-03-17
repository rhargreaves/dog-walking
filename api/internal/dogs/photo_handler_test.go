package dogs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDogRepository struct {
	mock.Mock
}

func (m *MockDogRepository) Get(id string) (*Dog, error) {
	args := m.Called(id)
	return args.Get(0).(*Dog), args.Error(1)
}

func (m *MockDogRepository) Create(dog *Dog) error {
	args := m.Called(dog)
	return args.Error(0)
}

func (m *MockDogRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDogRepository) List() ([]Dog, error) {
	args := m.Called()
	return args.Get(0).([]Dog), args.Error(1)
}

func (m *MockDogRepository) Update(id string, dog *Dog) error {
	args := m.Called(id, dog)
	return args.Error(0)
}

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
	const dogId = "123"

	dogRepository := new(MockDogRepository)
	dogRepository.On("Get", dogId).Return(&Dog{ID: dogId}, nil)
	handler := NewDogPhotoHandler(dogRepository, nil, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: dogId}}

	c.Request = httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/dogs/%s/photo", dogId), strings.NewReader("not an image"))
	c.Request.Header.Set("Content-Type", "text/plain")

	handler.UploadDogPhoto(c)

	requireAPIError(t, c, http.StatusBadRequest, "Invalid image content type")
}
