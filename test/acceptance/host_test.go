package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func skipIfLocal(t *testing.T) {
	if strings.HasPrefix(baseUrl(), "http://sam:") {
		t.Skip("Skipping TLS test on local environment")
	}
}

func TestServerIsRunning(t *testing.T) {
	resp := response(t, "/ping")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}

func TestServerIsNotAccessibleOnPort80(t *testing.T) {
	skipIfLocal(t)

	insecureUrl := strings.Replace(baseUrl(), "https://", "http://", 1)
	_, err := http.Get(insecureUrl)

	require.Error(t, err, "Expected connection to be refused on port 80")
	require.Contains(t, err.Error(), "connection refused",
		"Error should indicate connection refused")
}
