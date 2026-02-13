package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-ug/conveyor/internal/config"
	"github.com/open-ug/conveyor/internal/config/initialize"
	"github.com/open-ug/conveyor/pkg/server"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/stretchr/testify/assert"
)

// Build example resource definition (based on your example)
var payload = types.ResourceDefinition{
	Name:        "pipe3",
	Description: "Pipeline resource definition",
	Version:     "1.0.0",
	Schema: map[string]interface{}{
		"properties": map[string]interface{}{
			"image": map[string]interface{}{
				"type": "string",
			},
			"steps": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"command": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"name", "command"},
				},
			},
		},
		"required": []string{"steps"},
	},
}

func Test_ResourceDefinition_CRUD(t *testing.T) {
	configFile, err := initialize.Run(&initialize.Options{
		Force:   true,
		TempDir: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}
	config.LoadTestEnvConfig(configFile)

	cfg, err := config.GetTestConfig()
	if err != nil {
		t.Fatalf("failed to get test config: %v", err)
	}

	appctx, err := server.Setup(&cfg)
	if err != nil {
		t.Fatalf("failed to setup api: %v", err)
	}

	app := appctx.App

	// --- Create ---
	t.Run("create", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/resource-definitions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("create request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println(string(respBody))
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created on create")

	})

	// --- Get ---
	t.Run("get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource-definitions/"+payload.Name, nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("get request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK on get")

		var got types.ResourceDefinition
		if assert.NoError(t, json.Unmarshal(respBody, &got), "unmarshal get response") {
			assert.Equal(t, payload.Name, got.Name)
			assert.Equal(t, payload.Version, got.Version)
		}
	})

	// --- Update ---
	t.Run("update", func(t *testing.T) {
		// modify description
		updatedPayload := payload
		updatedPayload.Description = "Updated pipeline description"
		updateBody, _ := json.Marshal(updatedPayload)
		req := httptest.NewRequest(http.MethodPut, "/resource-definitions/"+payload.Name, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("update request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK on update")

		var updated types.ResourceDefinition
		if assert.NoError(t, json.Unmarshal(respBody, &updated), "unmarshal update response") {
			assert.Equal(t, "Updated pipeline description", updated.Description)
		}
	})

	// --- Apply (create or update) ---
	t.Run("apply", func(t *testing.T) {
		// change version and call apply endpoint
		appliedPayload := payload
		appliedPayload.Version = "1.0.1"
		applyBody, _ := json.Marshal(appliedPayload)
		req := httptest.NewRequest(http.MethodPost, "/resource-definitions/apply", bytes.NewReader(applyBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("apply request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created on apply")

		var applied types.ResourceDefinition
		if assert.NoError(t, json.Unmarshal(respBody, &applied), "unmarshal apply response") {
			assert.Equal(t, "1.0.1", applied.Version)
		}
	})

	// --- Delete ---
	t.Run("delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/resource-definitions/"+payload.Name, nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("delete request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "expected 204 No Content on delete")
	})

	// --- Get after delete (should fail based on current handler behavior) ---
	t.Run("get-after-delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource-definitions/"+payload.Name, nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("get-after-delete request failed: %v", err)
		}
		defer resp.Body.Close()

		// Current implementation returns 500 when FindOne fails / not found.
		// If you change handler behavior to return 404, update this assertion accordingly.
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "expected 500 after deleted resource (handler's current behavior)")
	})
	appctx.ShutDown()
}
