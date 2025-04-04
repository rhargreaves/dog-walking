package dogs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/model"
	"github.com/stretchr/testify/require"
)

func setupRouter(handler DogHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(common.ErrorMiddleware)
	router.GET("/dogs", handler.ListDogs)
	return router
}

func TestListDogs_ReturnsMaxDogsByDefault(t *testing.T) {
	numberOfDogs := 25
	dogRepository := NewFakeDogRepository()
	for i := 0; i < numberOfDogs; i++ {
		dogRepository.Create(domain.Dog{
			Name: fmt.Sprintf("Dog %d", i),
		})
	}

	config := DogHandlerConfig{
		ImagesCdnBaseUrl: "https://example.com",
	}
	handler := NewDogHandler(config, dogRepository)
	router := setupRouter(handler)
	req, _ := http.NewRequest("GET", "/dogs", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var dogs domain.DogList
	err := json.Unmarshal(resp.Body.Bytes(), &dogs)
	require.NoError(t, err)

	require.Equal(t, numberOfDogs, len(dogs.Dogs), "Expected %d dogs to be returned", numberOfDogs)
}

func TestListDogs_ReturnsErrorWhenLimitTooHigh(t *testing.T) {
	dogRepository := NewFakeDogRepository()
	config := DogHandlerConfig{
		ImagesCdnBaseUrl: "https://example.com",
	}
	handler := NewDogHandler(config, dogRepository)

	router := setupRouter(handler)
	req, _ := http.NewRequest("GET", "/dogs?limit=26", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	var apiError common.APIErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &apiError)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, apiError.Error.Code)
	require.Equal(t, "Key: 'DogListQuery.Limit' Error:Field validation for 'Limit' failed on the 'max' tag", apiError.Error.Message)
}

func TestListDogs_ReturnsErrorWhenLimitTooLow(t *testing.T) {
	dogRepository := NewFakeDogRepository()
	config := DogHandlerConfig{
		ImagesCdnBaseUrl: "https://example.com",
	}
	handler := NewDogHandler(config, dogRepository)

	router := setupRouter(handler)
	req, _ := http.NewRequest("GET", "/dogs?limit=0", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	var apiError common.APIErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &apiError)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, apiError.Error.Code)
	require.Equal(t, "Key: 'DogListQuery.Limit' Error:Field validation for 'Limit' failed on the 'min' tag", apiError.Error.Message)
}

func TestListDogs_ReturnsPhotoUrlForDogsWithPhoto(t *testing.T) {
	dogRepository := NewFakeDogRepository()
	dog1, _ := dogRepository.Create(domain.Dog{
		Name:      "Dog 1",
		PhotoHash: "1234567890",
	})
	dog2, _ := dogRepository.Create(domain.Dog{
		Name:      "Dog 2",
		PhotoHash: "0987654321",
	})
	config := DogHandlerConfig{
		ImagesCdnBaseUrl: "https://example.com",
	}
	handler := NewDogHandler(config, dogRepository)
	router := setupRouter(handler)
	req, _ := http.NewRequest("GET", "/dogs", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var dogs model.DogListResponse
	err := json.Unmarshal(resp.Body.Bytes(), &dogs)
	require.NoError(t, err)
	require.Equal(t, 2, len(dogs.Dogs))
	require.Equal(t, "https://example.com/"+dog1.ID+"?h=1234567890", dogs.Dogs[0].PhotoUrl)
	require.Equal(t, "https://example.com/"+dog2.ID+"?h=0987654321", dogs.Dogs[1].PhotoUrl)
}
