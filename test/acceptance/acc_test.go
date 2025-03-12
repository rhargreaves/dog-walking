package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func url() string {
	host := os.Getenv("API_HOST")
	return "https://" + host + "/hello"
}

func response(t *testing.T) *http.Response {
	resp, err := http.Get(url())
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}

func TestServerIsRunningOverTLS(t *testing.T) {
	resp := response(t)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}

func TestServerIsNotAccessibleOnPort80(t *testing.T) {
	_, err := http.Get("http://" + os.Getenv("API_HOST"))

	require.Error(t, err, "Expected connection to be refused on port 80")
	assert.Contains(t, err.Error(), "connection refused",
		"Error should indicate connection refused")
}
