package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-ug/conveyor/cmd/api"
	"github.com/open-ug/conveyor/internal/config"
	"github.com/open-ug/conveyor/pkg/types"
)

// Build example resource definition (based on your example)
var resource_definition = types.ResourceDefinition{
	Name:        "pipe4",
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

type PipelineResource struct {
	Name     string               `json:"name"`
	Resource string               `json:"resource"`
	Spec     PipelineResourceSpec `json:"spec"`
}

type PipelineResourceSpec struct {
	Image string         `json:"image"`
	Steps []PipelineStep `json:"steps"`
}

type PipelineStep struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

var resource = PipelineResource{
	Name:     "pipeline-1",
	Resource: "pipe4",
	Spec: PipelineResourceSpec{
		Image: "ubuntu:latest",
		Steps: []PipelineStep{
			{
				Name:    "list dir",
				Command: "ls",
			},
		},
	},
}

func Test_Resource_CRUD(t *testing.T) {
	config.InitConfig()

	// setup app (assumes api.Setup configures routes and dependencies for tests)
	appctx, err := api.Setup()
	if err != nil {
		t.Fatalf("failed to setup api: %v", err)
	}

	app := appctx.App

	// --- Create Resource Definition ---
	t.Run("create-resource-definition", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(resource_definition)
		req := httptest.NewRequest(http.MethodPost, "/resource-definitions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("create resource-definition request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created on create resource-definition")

		var created types.ResourceDefinition
		if assert.NoError(t, json.Unmarshal(respBody, &created), "unmarshal create resource-definition response") {
			assert.Equal(t, resource_definition.Name, created.Name)
			assert.Equal(t, resource_definition.Version, created.Version)
		}
	})

	// --- Create Resource ---
	t.Run("create-resource", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(resource)
		req := httptest.NewRequest(http.MethodPost, "/resources", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("create resource request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created on create resource")

		// response is a map with keys: name, runid, message
		var respMap map[string]interface{}
		if assert.NoError(t, json.Unmarshal(respBody, &respMap), "unmarshal create resource response") {
			assert.Equal(t, resource.Name, respMap["name"])
			assert.Contains(t, respMap["message"], "Resource created", "message should indicate creation")
			assert.NotEmpty(t, respMap["runid"], "runid should be present")
		}
	})

	// --- Get Resource ---
	t.Run("get-resource", func(t *testing.T) {
		url := "/resources/" + resource.Resource + "/" + resource.Name
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("get resource request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK on get resource")

		var got PipelineResource
		if assert.NoError(t, json.Unmarshal(respBody, &got), "unmarshal get resource response") {
			assert.Equal(t, resource.Name, got.Name)
			assert.Equal(t, resource.Resource, got.Resource)
			assert.Equal(t, resource.Spec.Image, got.Spec.Image)
			if assert.Len(t, got.Spec.Steps, 1) {
				assert.Equal(t, resource.Spec.Steps[0].Command, got.Spec.Steps[0].Command)
			}
		}
	})

	// --- Update Resource ---
	t.Run("update-resource", func(t *testing.T) {
		updated := resource
		updated.Spec.Image = "alpine:latest"

		bodyBytes, _ := json.Marshal(updated)
		url := "/resources/" + resource.Resource + "/" + resource.Name
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("update resource request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK on update resource")

		var got PipelineResource
		if assert.NoError(t, json.Unmarshal(respBody, &got), "unmarshal update response") {
			assert.Equal(t, "alpine:latest", got.Spec.Image)
		}
	})

	// --- Delete Resource ---
	t.Run("delete-resource", func(t *testing.T) {
		url := "/resources/" + resource.Resource + "/" + resource.Name
		req := httptest.NewRequest(http.MethodDelete, url, nil)

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("delete resource request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "expected 204 No Content on delete resource")
	})

	// --- Get after delete (should fail based on current handler behavior) ---
	t.Run("get-after-delete", func(t *testing.T) {
		url := "/resources/" + resource.Resource + "/" + resource.Name
		req := httptest.NewRequest(http.MethodGet, url, nil)

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
