package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jwtToken string

func GetJwtToken() string {
	if jwtToken == "" {
		if IsLocal() {
			fmt.Println("🔑 Authenticating with local credentials")
			jwtToken = auth.CreateLocalJWT()
		} else {
			fmt.Println("🔑 Authenticating with AWS Cognito")
			jwtToken = auth.GetCognitoJWT()
		}
	}
	return jwtToken
}

func IsLocal() bool {
	return os.Getenv("USE_LOCALSTACK") == "true"
}

func BaseUrl() string {
	apiBaseUrl := os.Getenv("API_BASE_URL")
	if apiBaseUrl == "" {
		log.Fatal("API_BASE_URL environment variable is required")
	}
	return apiBaseUrl
}

func ApiRequest(t *testing.T, method, endpoint string, auth bool, body io.Reader) *http.Request {
	var req *http.Request
	if auth {
		req = authedRequest(t, method, endpoint, body)
	} else {
		req = request(t, method, endpoint, body)
	}
	return req
}

func request(t *testing.T, method, endpoint string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, BaseUrl()+endpoint, body)
	require.NoError(t, err, "failed to create "+method+" request")
	return req
}

func authedRequest(t *testing.T, method, endpoint string, body io.Reader) *http.Request {
	req := request(t, method, endpoint, body)
	req.Header.Set("Authorization", "Bearer "+GetJwtToken())
	return req
}

func GetResponse(t *testing.T, req *http.Request) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to fetch URL")
	return resp
}

func Get(t *testing.T, endpoint string, auth bool) *http.Response {
	req := ApiRequest(t, http.MethodGet, endpoint, auth, nil)
	return GetResponse(t, req)
}

func PostJson(t *testing.T, endpoint string, body any, auth bool) *http.Response {
	req := ApiRequest(t, http.MethodPost, endpoint, auth, structToJsonBuffer(t, body))
	req.Header.Set("Content-Type", "application/json")
	return GetResponse(t, req)
}

func PutJson(t *testing.T, endpoint string, body any, auth bool) *http.Response {
	req := ApiRequest(t, http.MethodPut, endpoint, auth, structToJsonBuffer(t, body))
	req.Header.Set("Content-Type", "application/json")
	return GetResponse(t, req)
}

func Delete(t *testing.T, endpoint string, auth bool) *http.Response {
	req := ApiRequest(t, http.MethodDelete, endpoint, auth, nil)
	return GetResponse(t, req)
}

func CorsPreflight(t *testing.T, endpoint string, origin string, method string) *http.Response {
	req, err := http.NewRequest(http.MethodOptions, BaseUrl()+endpoint, nil)
	require.NoError(t, err, "failed to create OPTIONS request")
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", method)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")
	return GetResponse(t, req)
}

func RequireStatus(t *testing.T, resp *http.Response, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"expected status code %d", expectedStatus)

	if resp.StatusCode != expectedStatus {
		logBody(t, resp)
	}
}

func DecodeJSON(t *testing.T, resp *http.Response, target any) {
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
	if IsLocal() {
		t.Skip("Skipping test on local environment")
	}
}

func structToJsonBuffer(t *testing.T, body any) *bytes.Buffer {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "failed to marshal body")
	return bytes.NewBuffer(jsonBody)
}
