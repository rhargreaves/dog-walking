package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	jwtToken string
)

func isLocal() bool {
	return os.Getenv("USE_LOCALSTACK") == "true"
}

func Authenticate() {
	if isLocal() {
		fmt.Println("🔑 Authenticating with local credentials")
		jwtToken = auth.CreateLocalJWT()
	} else {
		fmt.Println("🔑 Authenticating with AWS Cognito")
		jwtToken = auth.GetCognitoJWT()
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func BaseUrl() string {
	return os.Getenv("API_BASE_URL")
}

func NewAuthedRequest(t *testing.T, method, endpoint string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, BaseUrl()+endpoint, body)
	require.NoError(t, err, "failed to create "+method+" request")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	return req
}

func Get(t *testing.T, endpoint string) *http.Response {
	req := NewAuthedRequest(t, http.MethodGet, endpoint, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to fetch URL")
	return resp
}

func PostBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	req := NewAuthedRequest(t, http.MethodPost, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err, "failed to POST to URL")
	return resp
}

func PostJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "failed to marshal body")
	return PostBytes(t, endpoint, jsonBody)
}

func PutBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	req := NewAuthedRequest(t, http.MethodPut, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to PUT to URL")
	return resp
}

func PutJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "failed to marshal body")
	return PutBytes(t, endpoint, jsonBody)
}

func Delete(t *testing.T, endpoint string) *http.Response {
	req := NewAuthedRequest(t, http.MethodDelete, endpoint, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to DELETE to URL")
	return resp
}

func RequireStatus(t *testing.T, resp *http.Response, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"expected status code %d", expectedStatus)

	if resp.StatusCode != expectedStatus {
		logBody(t, resp)
	}
}

func DecodeJSON(t *testing.T, resp *http.Response, target interface{}) {
	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "failed to decode response body")
}

func logBody(t *testing.T, resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	t.Errorf("Body: %s", string(bodyBytes))
}

func SkipIfLocal(t *testing.T) {
	if isLocal() {
		t.Skip("Skipping test on local environment")
	}
}
