package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func baseUrl() string {
	return os.Getenv("API_BASE_URL")
}

func get(t *testing.T, endpoint string) *http.Response {
	resp, err := http.Get(baseUrl() + endpoint)
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}

func postBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	resp, err := http.Post(baseUrl()+endpoint, "application/json", bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to POST to URL")
	return resp
}

func postJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal body")
	return postBytes(t, endpoint, jsonBody)
}
