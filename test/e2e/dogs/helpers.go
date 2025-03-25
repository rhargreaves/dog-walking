package dogs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func createDog(t *testing.T, name string) DogResponse {
	resp := common.PostJson(t, "/dogs", DogResponse{Name: name}, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusCreated)

	var dog DogResponse
	common.DecodeJSON(t, resp, &dog)

	assert.Equal(t, name, dog.Name, "Expected dog name to be returned")
	assert.NotEmpty(t, dog.ID, "Expected dog ID to be returned")

	return dog
}

func deleteDog(t *testing.T, id string) {
	resp := common.Delete(t, "/dogs/"+id, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusNoContent)
}

func FindFirst[T any](items []T, predicate func(T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}
