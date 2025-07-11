package driverruntime

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	t "github.com/open-ug/conveyor/pkg/types"
	"github.com/spf13/viper"
)

/*
* Broad Cast An event message over nats Network.
 */
func BroadCastMessage(
	message t.DriverMessage,
) error {
	fmt.Println("BroadCasting Message")
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := "http://" + viper.GetString("api.host") + ":" + viper.GetString("api.port")
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
