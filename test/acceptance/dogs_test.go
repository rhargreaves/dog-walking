package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Dog struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func logResponseBody(t *testing.T, resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	t.Errorf("Body is: %s", string(bodyBytes))
}

func TestDogsCreateDog(t *testing.T) {
	resp := postJson(t, "/dogs", Dog{Name: "Rover"})
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	if resp.StatusCode >= http.StatusInternalServerError {
		logResponseBody(t, resp)
	}

	var dogs Dog
	err := json.NewDecoder(resp.Body).Decode(&dogs)
	require.NoError(t, err, "Failed to decode response body")

	assert.Equal(t, "Rover", dogs.Name, "Expected dog name to be returned")
	assert.NotEmpty(t, dogs.ID, "Expected dog ID to be returned")
}

func TestDogsListDogs(t *testing.T) {
	resp := get(t, "/dogs")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	if resp.StatusCode >= http.StatusInternalServerError {
		logResponseBody(t, resp)
	}

	var dogs []Dog
	err := json.NewDecoder(resp.Body).Decode(&dogs)
	require.NoError(t, err, "Failed to decode response body")

	require.GreaterOrEqual(t, len(dogs), 1, "Expected at least 1 dog to be returned")
	require.NotEmpty(t, dogs[0].ID, "Expected dog ID to be returned")
}

func TestDogsRejectsInvalidJson(t *testing.T) {
	resp := postBytes(t, "/dogs", []byte("foo"))
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code 403")
}
