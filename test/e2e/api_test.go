package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func TestApi_Running(t *testing.T) {
	resp := common.Get(t, "/ping")
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusOK)
}

func TestApi_NotOnPort80(t *testing.T) {
	skipIfLocal(t)

	insecureUrl := strings.Replace(common.BaseUrl(), "https://", "http://", 1)
	_, err := http.Get(insecureUrl)

	require.Error(t, err, "Expected connection to be refused on port 80")
	require.Contains(t, err.Error(), "connection refused",
		"Error should indicate connection refused")
}

func skipIfLocal(t *testing.T) {
	if strings.HasPrefix(common.BaseUrl(), "http://sam:") {
		t.Skip("Skipping TLS test on local environment")
	}
}
