package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func BaseUrl() string {
	return os.Getenv("API_BASE_URL")
}

func Get(t *testing.T, endpoint string) *http.Response {
	resp, err := http.Get(BaseUrl() + endpoint)
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}

func PostBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	resp, err := http.Post(BaseUrl()+endpoint, "application/json", bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to POST to URL")
	return resp
}

func PostJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal body")
	return PostBytes(t, endpoint, jsonBody)
}

func PutBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	req, err := http.NewRequest(http.MethodPut, BaseUrl()+endpoint, bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to create PUT request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to PUT to URL")
	return resp
}

func PutJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal body")
	return PutBytes(t, endpoint, jsonBody)
}

func Delete(t *testing.T, endpoint string) *http.Response {
	req, err := http.NewRequest(http.MethodDelete, BaseUrl()+endpoint, nil)
	require.NoError(t, err, "Failed to create DELETE request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to DELETE to URL")
	return resp
}

func RequireStatus(t *testing.T, resp *http.Response, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"Expected status code %d", expectedStatus)

	if resp.StatusCode != expectedStatus {
		logBody(t, resp)
	}
}

func DecodeJSON(t *testing.T, resp *http.Response, target interface{}) {
	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "Failed to decode response body")
}

func logBody(t *testing.T, resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	t.Errorf("Body: %s", string(bodyBytes))
}

func SkipIfLocal(t *testing.T) {
	if strings.HasPrefix(BaseUrl(), "http://sam:") {
		t.Skip("Skipping test on local environment")
	}
}
