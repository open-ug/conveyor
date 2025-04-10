package client

import (
	"context"
	"encoding/json"
	"fmt"

	craneTypes "conveyor.cloud.cranom.tech/pkg/types"
	"github.com/go-resty/resty/v2"
)

/*
	`CreateApplication` creates a new application in the Conveyor API.

It takes a context and an application object as parameters.
*/
func (c *Client) CreateApplication(ctx context.Context, app *craneTypes.Application) (*craneTypes.Application, error) {
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
		Post(baseURL + "/applications/")
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	var responseApp craneTypes.Application
	err = json.Unmarshal(resp.Body(), &responseApp)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &responseApp, nil
}
