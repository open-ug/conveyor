package driverruntime

import (
	"encoding/json"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func BroadCastMessage(
	message craneTypes.DriverMessage,
) error {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := viper.GetString("api.host") + ":" + viper.GetString("api.port")
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	resp, err := client.R().
		SetBody(jsonMessage).
		Post(baseURL + "/drivers/broadcast-message")
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	fmt.Println("Response: ", resp)
	return nil
}
