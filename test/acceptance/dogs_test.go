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

func TestDogsCreateDog(t *testing.T) {
	resp := postJson(t, "/dogs", Dog{Name: "Rover"})
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	var dogs Dog
	err := json.NewDecoder(resp.Body).Decode(&dogs)
	require.NoError(t, err, "Failed to decode response body")

	require.Equal(t, "Rover", dogs.Name, "Expected dog name to be returned")
}

func TestDogsListDogs(t *testing.T) {
	resp := get(t, "/dogs")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	var dogs []Dog
	err := json.NewDecoder(resp.Body).Decode(&dogs)
	require.NoError(t, err, "Failed to decode response body")

	require.Equal(t, 1, len(dogs), "Expected 1 dog to be returned")
	require.Equal(t, "Rover", dogs[0].Name, "Expected dog name to be returned")
}

func TestDogsRejectsInvalidJson(t *testing.T) {
	resp := postBytes(t, "/dogs", []byte(`foo`))
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code 403")
}
