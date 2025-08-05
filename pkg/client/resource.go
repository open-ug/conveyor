package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	types "github.com/open-ug/conveyor/pkg/types"
)

// doRequest is a generic helper function that handles HTTP requests with JSON marshaling/unmarshaling
// This eliminates repetitive code across all API methods
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	req := client.R().SetContext(ctx)

	// Set request body if provided (for POST/PUT requests)
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("Error marshaling request body: %v\n", err)
			return fmt.Errorf("doRequest: failed to marshal request body: %w", err)
		}
		req.SetBody(jsonBody)
	}

	// Execute the request
	baseURL := c.GetAPIURL()
	resp, err := req.Execute(method, baseURL+path)
	if err != nil {
		fmt.Printf("Error making %s request to %s: %v\n", method, path, err)
		return fmt.Errorf("doRequest: request failed to execute: %w", err)
	}

	// Check for HTTP error status codes
	if resp.IsError() {
		fmt.Printf("Server returned error %d for %s %s: %s\n", resp.StatusCode(), method, path, string(resp.Body()))
		return fmt.Errorf("doRequest: server returned error %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Unmarshal response if result pointer is provided
	if result != nil {
		if err := json.Unmarshal(resp.Body(), result); err != nil {
			fmt.Printf("Error unmarshaling response: %v\n", err)
			return fmt.Errorf("doRequest: failed to unmarshal response: %w", err)
		}
	}

	return nil
}

/*
Creates a new resource in the Conveyor API.
This function is used to create a new resource, which can be any type of resource defined in the Conveyor API.
It takes a context and a Resource object as input, and returns an APIResponse object or an error.
It is important to ensure that the Resource object is properly defined according to the API specifications.
*/
func (c *Client) CreateResource(ctx context.Context, app *types.Resource) (*types.APIResponse, error) {
	var response types.APIResponse
	if err := c.doRequest(ctx, http.MethodPost, "/resources/", app, &response); err != nil {
		return nil, fmt.Errorf("CreateResource: failed to create resource %w", err)
	}
	return &response, nil
}

/*
Creates a new Resource Definition in the Conveyor API.
This is used to define the structure of a resource.
It allows for the creation of custom resource types.
It is typically used to define the schema for resources that can be created and managed by the Conveyor API.
*/
func (c *Client) CreateResourceDefinition(ctx context.Context, app *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	var response types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodPost, "/resource-definitions/", app, &response); err != nil {
		return nil, fmt.Errorf("CreateResourceDefinition: failed to create resource definition %w", err)
	}
	return &response, nil
}

/*
Creates or updates a Resource Definition in the Conveyor API.
This function is used to apply a Resource Definition, which can either create a new one or update an existing one.
This is useful for managing resource definitions in a declarative manner.
*/
func (c *Client) CreateOrUpdateResourceDefinition(ctx context.Context, app *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	var response types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodPost, "/resource-definitions/apply", app, &response); err != nil {
		return nil, fmt.Errorf("CreateOrUpdateResourceDefinition: failed to create or update resource definition %w", err)
	}
	return &response, nil
}

/*
Gets a Resource by its name from the Conveyor API.
This function retrieves a specific resource by its name.
It is useful for fetching details of a resource that has been previously created.
*/
func (c *Client) GetResource(ctx context.Context, name string, resourceDefinition string) (*types.Resource, error) {
	var response types.Resource
	path := fmt.Sprintf("/resources/%s/%s", resourceDefinition, name)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, fmt.Errorf("GetResource: failed to get resource %w", err)
	}
	return &response, nil
}

/*
Gets a Resource Definition by its name from the Conveyor API.
This function retrieves a specific resource definition by its name.
It is useful for fetching the schema or structure of a resource that has been previously defined.
*/
func (c *Client) GetResourceDefinition(ctx context.Context, name string) (*types.ResourceDefinition, error) {
	var response types.ResourceDefinition
	path := fmt.Sprintf("/resource-definitions/%s", name)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, fmt.Errorf("GetResourceDefinition: failed to get resource definition %w", err)
	}
	return &response, nil
}

/*
Updates a Resource by its name in the Conveyor API.
This function updates an existing resource with new data.
It is useful for modifying the properties of a resource that has been previously created.
*/
func (c *Client) UpdateResource(ctx context.Context, resource *types.Resource) (*types.Resource, error) {
	var response types.Resource
	path := fmt.Sprintf("/resources/%s/%s", resource.Resource, resource.Name)
	if err := c.doRequest(ctx, http.MethodPut, path, resource, &response); err != nil {
		return nil, fmt.Errorf("UpdateResource: failed to update resource %w", err)
	}
	return &response, nil
}

/*
Updates a Resource Definition by its name in the Conveyor API.
This function updates an existing resource definition with new data.
It is useful for modifying the schema or structure of a resource that has been previously defined.
*/
func (c *Client) UpdateResourceDefinition(ctx context.Context, resourceDefinition *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	var response types.ResourceDefinition
	path := fmt.Sprintf("/resource-definitions/%s", resourceDefinition.Name)
	if err := c.doRequest(ctx, http.MethodPut, path, resourceDefinition, &response); err != nil {
		return nil, fmt.Errorf("UpdateResourceDefinition: failed to update resource definition %w", err)
	}
	return &response, nil
}

/*
Deletes a Resource by its name in the Conveyor API.
This function deletes an existing resource.
It is useful for removing a resource that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResource(ctx context.Context, name string, resourceDefinition string) (*types.APIResponse, error) {
	var response types.APIResponse
	path := fmt.Sprintf("/resources/%s/%s", resourceDefinition, name)
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &response); err != nil {
		return nil, fmt.Errorf("DeleteResource: failed to delete resource %w", err)
	}
	return &response, nil
}

/*
Deletes a Resource Definition by its name in the Conveyor API.
This function deletes an existing resource definition.
It is useful for removing a resource definition that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResourceDefinition(ctx context.Context, name string) (*types.APIResponse, error) {
	var response types.APIResponse
	path := fmt.Sprintf("/resource-definitions/%s", name)
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &response); err != nil {
		return nil, fmt.Errorf("DeleteResourceDefinition: failed to delete resource definition %w", err)
	}
	return &response, nil
}
