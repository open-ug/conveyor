package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/open-ug/conveyor/pkg/types"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateResource(resource types.Resource, definition types.ResourceDefinition) (bool, error) {
	// Marshal the resource struct to JSON
	resourceJSON, err := json.Marshal(resource.Spec)
	if err != nil {
		return false, fmt.Errorf("failed to marshal resource: %w", err)
	}

	// Unmarshal the JSON into a generic interface{}
	var resourceData interface{}
	if err := json.Unmarshal(resourceJSON, &resourceData); err != nil {
		return false, fmt.Errorf("failed to unmarshal resource to interface{}: %w", err)
	}

	jsonSchema, err := json.Marshal(definition.Schema)
	if err != nil {
		return false, fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Compile the schema from the definition
	schema, err := jsonschema.CompileString("inline-schema", string(jsonSchema))
	if err != nil {
		return false, fmt.Errorf("failed to compile schema: %w", err)
	}

	// Validate the resource data
	if err := schema.Validate(resourceData); err != nil {
		return false, fmt.Errorf("validation failed: %w", err)
	}

	return true, nil

}
