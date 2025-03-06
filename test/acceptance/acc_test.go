package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getResponse(t *testing.T) *http.Response {
	host := os.Getenv("API_HOST")
	resp, err := http.Get("http://" + host)
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}

func TestServerIsRunning(t *testing.T) {
	resp := getResponse(t)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}
