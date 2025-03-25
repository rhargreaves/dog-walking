package dogs

import (
	"net/http"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func getDog(t *testing.T, id string) DogResponse {
	resp := common.Get(t, "/dogs/"+id, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
	var dog DogResponse
	common.DecodeJSON(t, resp, &dog)
	return dog
}

func createDog(t *testing.T, request CreateOrUpdateDogRequest) DogResponse {
	resp := common.PostJson(t, "/dogs", request, true)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusCreated)
	var dog DogResponse
	common.DecodeJSON(t, resp, &dog)
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
