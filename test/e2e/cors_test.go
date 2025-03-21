package main

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func allowedOrigin() string {
	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		log.Fatal("CORS_ALLOWED_ORIGIN is not set")
	}
	return allowedOrigin
}

func TestCors_PreflightRequest(t *testing.T) {
	allowedOrigin := allowedOrigin()
	endpoints := []struct {
		name   string
		path   string
		method string
	}{
		{"Ping Endpoint", "/ping", "GET"},
		{"Dogs List", "/dogs", "GET"},
		{"Create Dog", "/dogs", "POST"},
		{"Update Dog", "/dogs/123", "PUT"},
		{"Delete Dog", "/dogs/123", "DELETE"},
		{"Upload Photo", "/dogs/123/photo", "PUT"},
		{"Detect Breed", "/dogs/123/photo/detect-breed", "POST"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			resp := common.CorsPreflight(t, endpoint.path, allowedOrigin, endpoint.method)
			defer resp.Body.Close()

			statusOK := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent
			assert.True(t, statusOK, "Expected status code 200 or 204, got %d", resp.StatusCode)

			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			assert.Equal(t, allowedOrigin, allowOrigin, "Access-Control-Allow-Origin should match")

			allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
			assert.Contains(t, allowMethods, endpoint.method, "Access-Control-Allow-Methods should contain the requested method")

			allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
			assert.Contains(t, allowHeaders, "content-type", "Access-Control-Allow-Headers should contain content-type")
			assert.Contains(t, allowHeaders, "authorization", "Access-Control-Allow-Headers should contain authorization")

			allowCredentials := resp.Header.Get("Access-Control-Allow-Credentials")
			assert.Equal(t, "true", allowCredentials, "Access-Control-Allow-Credentials should be true")
		})
	}
}

func TestCors_DisallowedOrigin(t *testing.T) {
	disallowedOrigin := "https://example.com"

	resp := common.CorsPreflight(t, "/ping", disallowedOrigin, "GET")
	defer resp.Body.Close()

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.NotEqual(t, disallowedOrigin, allowOrigin, "Access-Control-Allow-Origin should not match disallowed origin")
}

func TestCors_NormalRequestWithAllowedOrigin(t *testing.T) {
	allowedOrigin := allowedOrigin()

	req, err := http.NewRequest(http.MethodGet, common.BaseUrl()+"/ping", nil)
	require.NoError(t, err, "failed to create GET request")
	req.Header.Set("Origin", allowedOrigin)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to make GET request with Origin header")
	defer resp.Body.Close()

	common.RequireStatus(t, resp, http.StatusOK)
	respAllowedOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.Equal(t, allowedOrigin, respAllowedOrigin, "Access-Control-Allow-Origin should match in normal response")

	allowCredentials := resp.Header.Get("Access-Control-Allow-Credentials")
	assert.Equal(t, "true", allowCredentials, "Access-Control-Allow-Credentials should be true in normal response")
}
