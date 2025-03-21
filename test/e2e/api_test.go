package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func TestApi_Running(t *testing.T) {
	resp := common.Get(t, "/ping", false)
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
}

func TestApi_NotOnPort80(t *testing.T) {
	common.SkipIfLocal(t)

	insecureUrl := strings.Replace(common.BaseUrl(), "https://", "http://", 1)
	_, err := http.Get(insecureUrl)

	require.Error(t, err, "Expected connection to be refused on port 80")
	require.Contains(t, err.Error(), "connection refused",
		"Error should indicate connection refused")
}

func TestApi_AuthRequiredOnProtectedRoutes(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		endpoint string
	}{
		{"POST /dogs", http.MethodPost, "/dogs"},
		{"GET /dogs", http.MethodGet, "/dogs"},
		{"GET /dogs/{id}", http.MethodGet, "/dogs/123"},
		{"PUT /dogs/{id}", http.MethodPut, "/dogs/123"},
		{"DELETE /dogs/{id}", http.MethodDelete, "/dogs/123"},
		{"PUT /dogs/{id}/photo", http.MethodPut, "/dogs/123/photo"},
		{"POST /dogs/{id}/photo/detect-breed", http.MethodPost, "/dogs/123/photo/detect-breed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			switch tc.method {
			case http.MethodGet:
				resp = common.Get(t, tc.endpoint, false)
			case http.MethodPost:
				resp = common.PostJson(t, tc.endpoint, struct{}{}, false)
			case http.MethodPut:
				resp = common.PutJson(t, tc.endpoint, struct{}{}, false)
			case http.MethodDelete:
				resp = common.Delete(t, tc.endpoint, false)
			}
			defer resp.Body.Close()
			common.RequireStatus(t, resp, http.StatusUnauthorized)

			var messageResponse common.MessageResponse
			common.DecodeJSON(t, resp, &messageResponse)
			assert.Equal(t, "Unauthorized", messageResponse.Message,
				"Expected error message to be returned")
		})
	}
}
