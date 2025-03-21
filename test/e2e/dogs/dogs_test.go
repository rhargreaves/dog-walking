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
	dog := createDog(t, "Rover")
	assert.Equal(t, "Rover", dog.Name)
	assert.NotEmpty(t, dog.ID)
}

func TestCreateDog_RejectsInvalidJson(t *testing.T) {
	req := common.ApiRequest(t, http.MethodPost, "/dogs", true, bytes.NewBuffer([]byte("foo")))
	req.Header.Set("Content-Type", "application/json")
	resp := common.GetResponse(t, req)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusBadRequest)
}

func TestGetDog_ReturnsDogWhenExists(t *testing.T) {
	const dogName = "Rover"
	createdDog := createDog(t, dogName)

	resp := common.Get(t, "/dogs/"+createdDog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog Dog
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, dogName, fetchedDog.Name, "Expected dog name to be returned")
	assert.Equal(t, createdDog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestGetDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.Get(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)

	var errorResponse common.ErrorResponse
	common.DecodeJSON(t, resp, &errorResponse)
	assert.Equal(t, "dog not found", errorResponse.Error,
		"Expected error message to be returned")
}

func TestUpdateDog_UpdatesDogWhenExists(t *testing.T) {
	dog := createDog(t, "Rover")

	resp := common.PutJson(t, "/dogs/"+dog.ID, Dog{Name: "Rose"}, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.Get(t, "/dogs/"+dog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog Dog
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, "Rose", fetchedDog.Name, "Expected updated dog name to be returned")
	assert.Equal(t, dog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestUpdateDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.PutJson(t, "/dogs/123", Dog{Name: "Mr. Peanutbutter"}, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)

	resp = common.Get(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}

func TestListDogs(t *testing.T) {
	createDog(t, "ListTest")

	resp := common.Get(t, "/dogs", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogs []Dog
	common.DecodeJSON(t, resp, &dogs)

	require.GreaterOrEqual(t, len(dogs), 1, "Expected at least 1 dog to be returned")
	require.NotEmpty(t, dogs[0].ID, "Expected dog ID to be returned")
}

func TestDeleteDog_DeletesDogWhenExists(t *testing.T) {
	dog := createDog(t, "Rover")

	resp := common.Delete(t, "/dogs/"+dog.ID, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
}

func TestDeleteDog_ReturnsNotFoundWhenDogDoesNotExist(t *testing.T) {
	resp := common.Delete(t, "/dogs/123", true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)
}
