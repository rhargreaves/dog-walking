package dogs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
)

func TestCreateDog(t *testing.T) {
	dog := createDog(t, "Rover")
	assert.Equal(t, "Rover", dog.Name)
	assert.NotEmpty(t, dog.ID)
}

func TestCreateDog_RejectsInvalidJson(t *testing.T) {
	resp := common.PostBytes(t, "/dogs", []byte("foo"))
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusBadRequest)
}

func TestGetDog_Exists(t *testing.T) {
	const dogName = "Rover"
	createdDog := createDog(t, dogName)

	resp := common.Get(t, "/dogs/"+createdDog.ID)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog Dog
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, dogName, fetchedDog.Name, "Expected dog name to be returned")
	assert.Equal(t, createdDog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestGetDog_ReturnsNotFoundForMissingDog(t *testing.T) {
	resp := common.Get(t, "/dogs/123")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNotFound)

	var errorResponse common.ErrorResponse
	common.DecodeJSON(t, resp, &errorResponse)
	assert.Equal(t, "dog not found", errorResponse.Error,
		"Expected error message to be returned")
}

func TestUpdateDog_Exists(t *testing.T) {
	dog := createDog(t, "Rover")

	resp := common.PutJson(t, "/dogs/"+dog.ID, Dog{Name: "Rose"})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	resp = common.Get(t, "/dogs/"+dog.ID)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var fetchedDog Dog
	common.DecodeJSON(t, resp, &fetchedDog)

	assert.Equal(t, "Rose", fetchedDog.Name, "Expected updated dog name to be returned")
	assert.Equal(t, dog.ID, fetchedDog.ID, "Expected dog ID to be returned")
}

func TestListDogs(t *testing.T) {
	createDog(t, "ListTest")

	resp := common.Get(t, "/dogs")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)

	var dogs []Dog
	common.DecodeJSON(t, resp, &dogs)

	require.GreaterOrEqual(t, len(dogs), 1, "Expected at least 1 dog to be returned")
	require.NotEmpty(t, dogs[0].ID, "Expected dog ID to be returned")
}
