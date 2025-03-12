package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func baseUrl() string {
	host := os.Getenv("API_HOST")
	return "https://" + host
}

func response(t *testing.T, endpoint string) *http.Response {
	resp, err := http.Get(baseUrl() + endpoint)
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}
