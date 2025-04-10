package webhookdriver

import (
	"encoding/json"
	"fmt"

	craneTypes "conveyor.cloud.cranom.tech/pkg/types"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func PostMessage(
	message craneTypes.DriverMessage,
) error {
	fmt.Println("BroadCasting Message")
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := viper.GetString("webhook.endpoint")
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		Post(baseURL)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	fmt.Println("Response: ", resp)
	return nil
}
