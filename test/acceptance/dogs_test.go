package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Dog struct {
	Name string `json:"name"`
}

func TestDogsReturnsEmptyListOfDogs(t *testing.T) {
	resp := postJson(t, "/dogs", Dog{})
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	var dogs Dog
	err := json.NewDecoder(resp.Body).Decode(&dogs)
	require.NoError(t, err, "Failed to decode response body")

	require.Empty(t, dogs, "Expected empty list of dogs")
}
