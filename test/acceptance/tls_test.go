package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerIsRunningOverTLS(t *testing.T) {
	resp := response(t, "/hello")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}

func TestServerIsNotAccessibleOnPort80(t *testing.T) {
	_, err := http.Get("http://" + os.Getenv("API_HOST"))

	require.Error(t, err, "Expected connection to be refused on port 80")
	require.Contains(t, err.Error(), "connection refused",
		"Error should indicate connection refused")
}
