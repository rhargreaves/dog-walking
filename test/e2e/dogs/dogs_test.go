package dogs

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func TestCreateDog(t *testing.T) {
	dog := createDog(t, testDog)
	assert.Equal(t, "Rover", dog.Name)
	assert.NotEmpty(t, dog.ID)

	fetchedDog := getDog(t, dog.ID)
	assert.Equal(t, "Rover", fetchedDog.Name)
	assert.Equal(t, "Labrador", fetchedDog.Breed)
	assert.Equal(t, "male", fetchedDog.Sex)
	assert.True(t, fetchedDog.IsNeutered)
	assert.Equal(t, 3, fetchedDog.EnergyLevel)
	assert.Equal(t, "medium", fetchedDog.Size)
	assert.Equal(t, "None", fetchedDog.SpecialInstructions)
	assert.Equal(t, "2020-01-01", fetchedDog.DateOfBirth)
}

func TestCreateDog_RejectsInvalidJson(t *testing.T) {
	req := common.ApiRequest(t, http.MethodPost, "/dogs", true, bytes.NewBuffer([]byte("foo")))
	req.Header.Set("Content-Type", "application/json")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusBadRequest)
}

func TestGetDog_ReturnsDogWhenExists(t *testing.T) {
	createdDog := createDog(t, testDog)

	resp := common.Get(t, "/dogs/"+createdDog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog DogResponse
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, testDog.Name, fetchedDog.Name, "Expected dog name to be returned")
	assert.Equal(t, createdDog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestGetDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.Get(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)

	var errorResponse common.ApiErrorResponse
	common.DecodeJSON(t, resp, &errorResponse)
	assert.Equal(t, "dog not found", errorResponse.Error.Message,
		"Expected error message to be returned")
}

func TestUpdateDog_UpdatesDogWhenExists(t *testing.T) {
	dog := createDog(t, testDog)

	resp := common.PutJson(t, "/dogs/"+dog.ID, testDog2, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.Get(t, "/dogs/"+dog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog DogResponse
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, "Echo", fetchedDog.Name, "Expected updated dog name to be returned")
	assert.Equal(t, "Husky", fetchedDog.Breed, "Expected updated dog breed to be returned")
	assert.Equal(t, "male", fetchedDog.Sex, "Expected updated dog sex to be returned")
	assert.False(t, fetchedDog.IsNeutered, "Expected updated dog isNeutered to be returned")
	assert.Equal(t, 5, fetchedDog.EnergyLevel, "Expected updated dog energyLevel to be returned")
	assert.Equal(t, "large", fetchedDog.Size, "Expected updated dog size to be returned")
	assert.Equal(t, "Don't let him out of the house", fetchedDog.SpecialInstructions, "Expected updated dog specialInstructions to be returned")
	assert.Equal(t, "2020-01-01", fetchedDog.DateOfBirth, "Expected updated dog dateOfBirth to be returned")
	assert.Equal(t, false, fetchedDog.Socialization.GoodWithChildren, "Expected updated dog socialization.goodWithChildren to be returned")
	assert.Equal(t, false, fetchedDog.Socialization.GoodWithPuppies, "Expected updated dog socialization.goodWithPuppies to be returned")
	assert.Equal(t, false, fetchedDog.Socialization.GoodWithLargeDogs, "Expected updated dog socialization.goodWithLargeDogs to be returned")
	assert.Equal(t, false, fetchedDog.Socialization.GoodWithSmallDogs, "Expected updated dog socialization.goodWithSmallDogs to be returned")
	assert.Equal(t, dog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestUpdateDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.PutJson(t, "/dogs/123", DogResponse{Name: "Mr. Peanutbutter"}, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)

	resp = common.Get(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}

func TestListDogs_ReturnsAtLeastOneDog(t *testing.T) {
	createDog(t, testListDog)

	resp := common.Get(t, "/dogs", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogs DogListResponse
	common.DecodeJSON(t, resp, &dogs)

	require.GreaterOrEqual(t, len(dogs.Dogs), 1, "Expected at least 1 dog to be returned")
	require.NotEmpty(t, dogs.Dogs[0].ID, "Expected dog ID to be returned")
}

func TestListDogs_ReturnsNextTokenWhenMoreDogsExist(t *testing.T) {
	createDog(t, testListDog)
	createDog(t, testListDog)

	resp := common.Get(t, "/dogs?limit=1", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogs DogListResponse
	common.DecodeJSON(t, resp, &dogs)

	require.Equal(t, 1, len(dogs.Dogs), "Expected 1 dog to be returned")
	require.NotEmpty(t, dogs.NextToken, "Expected next token to be returned")

	resp = common.Get(t, "/dogs?limit=1&nextToken="+dogs.NextToken, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogsNextPage DogListResponse
	common.DecodeJSON(t, resp, &dogsNextPage)
	require.NotEqual(t, dogs.Dogs[0].ID, dogsNextPage.Dogs[0].ID, "Expected next token to return a different dog")
}

func TestListDogs_ReturnsDogsFilteredByName(t *testing.T) {
	createDog(t, testNameFilterDog)
	createDog(t, testNameFilterDog)

	resp := common.Get(t, "/dogs?limit=2&name=NameFilterTest", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogs DogListResponse
	common.DecodeJSON(t, resp, &dogs)

	assert.Equal(t, 2, len(dogs.Dogs), "Expected 2 dogs to be returned")
	for _, dog := range dogs.Dogs {
		assert.Contains(t, dog.Name, "NameFilterTest", "Expected dog name to contain 'NameFilterTest'")
		deleteDog(t, dog.ID)
	}
}

func TestDeleteDog_DeletesDogWhenExists(t *testing.T) {
	dog := createDog(t, testDog)

	resp := common.Delete(t, "/dogs/"+dog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNoContent)
}

func TestDeleteDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.Delete(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}
