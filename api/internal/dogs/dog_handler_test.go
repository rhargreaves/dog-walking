package dogs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/models"
	"github.com/rhargreaves/dog-walking/api/internal/mocks"
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
	dogRepository := new(mocks.DogRepository)
	dogList := &models.DogList{
		Dogs:      make([]models.Dog, numberOfDogs),
		NextToken: "",
	}
	for i := range dogList.Dogs {
		dogList.Dogs[i] = models.Dog{
			ID:   fmt.Sprintf("dog-%d", i),
			Name: fmt.Sprintf("Dog %d", i),
		}
	}
	dogRepository.EXPECT().List(numberOfDogs, "").Return(dogList, nil)

	handler := NewDogHandler(dogRepository)
	router := setupRouter(handler)
	req, _ := http.NewRequest("GET", "/dogs", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var dogs models.DogList
	err := json.Unmarshal(resp.Body.Bytes(), &dogs)
	require.NoError(t, err)

	require.Equal(t, numberOfDogs, len(dogs.Dogs), "Expected %d dogs to be returned", numberOfDogs)
}

func TestListDogs_ReturnsErrorWhenLimitTooHigh(t *testing.T) {
	dogRepository := new(mocks.DogRepository)
	handler := NewDogHandler(dogRepository)

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
	dogRepository := new(mocks.DogRepository)
	handler := NewDogHandler(dogRepository)

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
