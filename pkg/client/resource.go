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
*/
func (c *Client) CreateResource(ctx context.Context, app *types.Resource) (*types.Resource, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := "http://" + c.APIHost + ":" + c.APIPort
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
	var responseApp types.Resource
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}

/*
Creates a new Resource Definition in the Conveyor API.
*/

func (c *Client) CreateResourceDefinition(ctx context.Context, app *types.ResourceDefinition) (*types.ResourceDefinition, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := "http://" + c.APIHost + ":" + c.APIPort
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
