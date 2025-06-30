package helpers_test

import (
	"testing"

	"github.com/open-ug/conveyor/internal/helpers"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateResource_ValidInput(t *testing.T) {
	resource := types.Resource{
		Spec: map[string]interface{}{
			"name":  "my-app",
			"image": "nginx:latest",
		},
	}

	definition := types.ResourceDefinition{
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
				"image": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"name", "image"},
		},
	}

	valid, err := helpers.ValidateResource(resource, definition)
	assert.True(t, valid)
	assert.NoError(t, err)
}

func TestValidateResource_InvalidInput(t *testing.T) {
	resource := types.Resource{
		Spec: map[string]interface{}{
			"name":  "my-app",
			"image": 123, // invalid type
		},
	}

	definition := types.ResourceDefinition{
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
				"image": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"name", "image"},
		},
	}

	valid, err := helpers.ValidateResource(resource, definition)
	assert.False(t, valid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}
