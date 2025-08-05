package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-ug/conveyor/pkg/types"
)

// doRequest is a helper function to make HTTP requests to the Conveyor API.
func (c *Client) doRequest(ctx context.Context, method, path string, body, dest any) error {
	if (method == http.MethodPost || method == http.MethodPut) && body == nil {
		return fmt.Errorf("doRequest: body cannot be nil for POST/PUT requests")
	}
	if dest == nil {
		return fmt.Errorf("doRequest: destination cannot be nil")
	}

	req := c.HTTPClient.R().SetContext(ctx)

	if body != nil {
		jsonMessage, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("doRequest: failed to marshal request body: %w", err)
		}
		req.SetBody(jsonMessage)
	}

	resp, err := req.Execute(method, c.HTTPClient.BaseURL+path)
	if err != nil {
		return fmt.Errorf("doRequest: failed to execute %s request: %w", method, err)
	}

	if err := json.Unmarshal(resp.Body(), dest); err != nil {
		return fmt.Errorf("doRequest: failed to unmarshal response body: %w", err)
	}

	return nil
}

/*
Creates a new resource in the Conveyor API.
This function is used to create a new resource, which can be any type of resource defined in the Conveyor API.
It takes a context and a Resource object as input, and returns an APIResponse object or an error.
It is important to ensure that the Resource object is properly defined according to the API specifications.
*/
func (c *Client) CreateResource(ctx context.Context, resource *types.Resource) (*types.APIResponse, error) {
	var resp types.APIResponse
	if err := c.doRequest(ctx, http.MethodPost, "/resources/", resource, &resp); err != nil {
		return nil, fmt.Errorf("CreateResource: failed to create resource, %w", err)
	}

	return &resp, nil
}

/*
Creates a new Resource Definition in the Conveyor API.
This is used to define the structure of a resource.
It allows for the creation of custom resource types.
It is typically used to define the schema for resources that can be created and managed by the Conveyor API.
*/
func (c *Client) CreateResourceDefinition(ctx context.Context, resourceDef *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	var resp types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodPost, "/resource-definitions/", resourceDef, &resp); err != nil {
		return nil, fmt.Errorf("CreateResourceDefinition: failed to create resource definition, %w", err)
	}

	return &resp, nil
}

/*
Creates or updates a Resource Definition in the Conveyor API.
This function is used to apply a Resource Definition, which can either create a new one or update an existing one.
This is useful for managing resource definitions in a declarative manner.
*/
func (c *Client) CreateOrUpdateResourceDefinition(ctx context.Context, resourceDef *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	var resp types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodPost, "/resource-definitions/apply", resourceDef, &resp); err != nil {
		return nil, fmt.Errorf("CreateOrUpdateResourceDefinition: failed to create or update resource definition, %w", err)
	}

	return &resp, nil
}

/*
Gets a Resource by its name from the Conveyor API.
This function retrieves a specific resource by its name.
It is useful for fetching details of a resource that has been previously created.
*/
func (c *Client) GetResource(ctx context.Context, name string, resourceDefinition string) (*types.Resource, error) {
	path := fmt.Sprintf("/resources/%s/%s", resourceDefinition, name)

	var resp types.Resource
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("GetResource: failed to get resource, %w", err)
	}

	return &resp, nil
}

/*
Gets a Resource Definition by its name from the Conveyor API.
This function retrieves a specific resource definition by its name.
It is useful for fetching the schema or structure of a resource that has been previously defined.
*/
func (c *Client) GetResourceDefinition(ctx context.Context, name string) (*types.ResourceDefinition, error) {
	path := fmt.Sprintf("/resource-definitions/%s", name)
	var resp types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("GetResourceDefinition: failed to get resource definition, %w", err)
	}

	return &resp, nil
}

/*
Updates a Resource by its name in the Conveyor API.
This function updates an existing resource with new data.
It is useful for modifying the properties of a resource that has been previously created.
*/
func (c *Client) UpdateResource(ctx context.Context, resource *types.Resource) (*types.Resource, error) {
	path := fmt.Sprintf("/resources/%s/%s", resource.Resource, resource.Name)

	var resp types.Resource
	if err := c.doRequest(ctx, http.MethodPut, path, resource, &resp); err != nil {
		return nil, fmt.Errorf("UpdateResource: failed to update resource, %w", err)
	}

	return &resp, nil
}

/*
Updates a Resource Definition by its name in the Conveyor API.
This function updates an existing resource definition with new data.
It is useful for modifying the schema or structure of a resource that has been previously defined.
*/
func (c *Client) UpdateResourceDefinition(ctx context.Context, resourceDefinition *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	path := fmt.Sprintf("/resource-definitions/%s", resourceDefinition.Name)

	var resp types.ResourceDefinition
	if err := c.doRequest(ctx, http.MethodPut, path, resourceDefinition, &resp); err != nil {
		return nil, fmt.Errorf("UpdateResourceDefinition: failed to update resource definition, %w", err)
	}

	return &resp, nil
}

/*
Deletes a Resource by its name in the Conveyor API.
This function deletes an existing resource.
It is useful for removing a resource that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResource(ctx context.Context, name string, resourceDefinition string) (*types.APIResponse, error) {
	path := fmt.Sprintf("/resources/%s/%s", resourceDefinition, name)

	var resp types.APIResponse
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("DeleteResource: failed to delete resource, %w", err)
	}

	return &resp, nil
}

/*
Deletes a Resource Definition by its name in the Conveyor API.
This function deletes an existing resource definition.
It is useful for removing a resource definition that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResourceDefinition(ctx context.Context, name string) (*types.APIResponse, error) {
	path := fmt.Sprintf("/resource-definitions/%s", name)

	var resp types.APIResponse
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("DeleteResourceDefinition: failed to delete resource definition, %w", err)
	}

	return &resp, nil
}
