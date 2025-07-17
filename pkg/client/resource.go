package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	types "github.com/open-ug/conveyor/pkg/types"
)

/*
Creates a new resource in the Conveyor API.
This function is used to create a new resource, which can be any type of resource defined in the Conveyor API.
It takes a context and a Resource object as input, and returns an APIResponse object or an error.
It is important to ensure that the Resource object is properly defined according to the API specifications.
*/
func (c *Client) CreateResource(ctx context.Context, app *types.Resource) (*types.APIResponse, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	jsonMessage, err := json.Marshal(app)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		SetContext(ctx).
		Post(baseURL + "/resources/")
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.APIResponse
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Creates a new Resource Definition in the Conveyor API.
This is used to define the structure of a resource.
It allows for the creation of custom resource types.
It is typically used to define the schema for resources that can be created and managed by the Conveyor API.
*/
func (c *Client) CreateResourceDefinition(ctx context.Context, app *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	jsonMessage, err := json.Marshal(app)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		SetContext(ctx).
		Post(baseURL + "/resource-definitions/")
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.ResourceDefinition
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Creates or updates a Resource Definition in the Conveyor API.
This function is used to apply a Resource Definition, which can either create a new one or update an existing one.
This is useful for managing resource definitions in a declarative manner.
*/
func (c *Client) CreateOrUpdateResourceDefinition(ctx context.Context, app *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	jsonMessage, err := json.Marshal(app)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		SetContext(ctx).
		Post(baseURL + "/resource-definitions/apply")
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.ResourceDefinition
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Gets a Resource by its name from the Conveyor API.
This function retrieves a specific resource by its name.
It is useful for fetching details of a resource that has been previously created.
*/
func (c *Client) GetResource(ctx context.Context, name string, resouce_definition string) (*types.Resource, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	resp, err := client.R().
		SetContext(ctx).
		Get(baseURL + "/resources/" + resouce_definition + "/" + name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.Resource
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Gets a Resource Definition by its name from the Conveyor API.
This function retrieves a specific resource definition by its name.
It is useful for fetching the schema or structure of a resource that has been previously defined.
*/

func (c *Client) GetResourceDefinition(ctx context.Context, name string) (*types.ResourceDefinition, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	resp, err := client.R().
		SetContext(ctx).
		Get(baseURL + "/resource-definitions/" + name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.ResourceDefinition
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Updates a Resource by its name in the Conveyor API.
This function updates an existing resource with new data.
It is useful for modifying the properties of a resource that has been previously created.
*/
func (c *Client) UpdateResource(ctx context.Context, resource *types.Resource) (*types.Resource, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	jsonMessage, err := json.Marshal(resource)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		SetContext(ctx).
		Put(baseURL + "/resources/" + resource.Resource + "/" + resource.Name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.Resource
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Updates a Resource Definition by its name in the Conveyor API.
This function updates an existing resource definition with new data.
It is useful for modifying the schema or structure of a resource that has been previously defined.
*/
func (c *Client) UpdateResourceDefinition(ctx context.Context, resourceDefinition *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	jsonMessage, err := json.Marshal(resourceDefinition)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		SetContext(ctx).
		Put(baseURL + "/resource-definitions/" + resourceDefinition.Name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.ResourceDefinition
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Deletes a Resource by its name in the Conveyor API.
This function deletes an existing resource.
It is useful for removing a resource that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResource(ctx context.Context, name string, resource_definition string) (*types.APIResponse, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	resp, err := client.R().
		SetContext(ctx).
		Delete(baseURL + "/resources/" + resource_definition + "/" + name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.APIResponse
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Deletes a Resource Definition by its name in the Conveyor API.
This function deletes an existing resource definition.
It is useful for removing a resource definition that is no longer needed or has been replaced.
*/
func (c *Client) DeleteResourceDefinition(ctx context.Context, name string) (*types.APIResponse, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := c.GetAPIURL()
	resp, err := client.R().
		SetContext(ctx).
		Delete(baseURL + "/resource-definitions/" + name)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp types.APIResponse
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}
