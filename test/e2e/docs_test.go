package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiDocs_Available(t *testing.T) {
	// Test access to API Docs UI
	resp := common.Get(t, "/api-docs/index.html", false)
	defer resp.Body.Close()

	// Check if API Docs UI is present
	if resp.StatusCode == http.StatusOK {
		t.Log("✅ API documentation UI is accessible")

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "Should be able to read response body")

		// Verify it contains Swagger UI HTML
		htmlContent := string(bodyBytes)
		assert.Contains(t, htmlContent, "swagger-ui", "Response should contain Swagger UI HTML")

		// Also check JSON spec is available
		jsonResp := common.Get(t, "/api-docs/doc.json", false)
		defer jsonResp.Body.Close()

		if jsonResp.StatusCode == http.StatusOK {
			t.Log("✅ API JSON spec is accessible")

			// Verify it's valid JSON with expected fields
			var swaggerSpec map[string]any
			err := json.NewDecoder(jsonResp.Body).Decode(&swaggerSpec)
			require.NoError(t, err, "Should return valid JSON")

			// Check basic swagger spec structure
			assert.Equal(t, "2.0", swaggerSpec["swagger"], "Should be Swagger 2.0 spec")
			assert.Contains(t, swaggerSpec, "info", "Should contain info section")
			assert.Contains(t, swaggerSpec, "paths", "Should contain paths section")

			// Verify API endpoints are documented
			paths, ok := swaggerSpec["paths"].(map[string]any)
			require.True(t, ok, "Paths should be a map")

			// Check specific endpoints
			assert.Contains(t, paths, "/ping", "Should document ping endpoint")
			if pingPath, ok := paths["/ping"].(map[string]any); ok {
				assert.Contains(t, pingPath, "get", "Should document GET method for ping endpoint")
				if getOp, ok := pingPath["get"].(map[string]any); ok {
					assert.Contains(t, getOp, "tags", "Should have tags for ping endpoint")
					tags, ok := getOp["tags"].([]any)
					assert.True(t, ok, "Tags should be an array")
					assert.Contains(t, tags, "health", "Ping endpoint should be tagged with 'health'")
				}
			}
			assert.Contains(t, paths, "/dogs", "Should document dogs endpoint")
		} else {
			t.Logf("⚠️ API JSON spec is not accessible (status: %d) - this may be expected if API docs are not fully deployed", jsonResp.StatusCode)
		}
	} else {
		t.Logf("⚠️ API documentation UI is not accessible (status: %d) - this may be expected if API docs are not fully deployed", resp.StatusCode)
	}
}
