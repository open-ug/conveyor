package dockerdriver

import (
	"context"
	"fmt"

	apiServer "crane.cloud.cranom.tech/cmd/api"
)

func Listen() {
	rdb := apiServer.NewRedisClient()

	pubsub := rdb.Subscribe(context.Background(), "application")

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}

}
