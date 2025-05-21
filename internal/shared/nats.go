/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package shared

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func NewNatsConn() *nats.Conn {
	nc, err := nats.Connect(viper.GetString("nats.url"))
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	return nc
}
