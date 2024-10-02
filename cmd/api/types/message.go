package types

type DriverMessage struct {
	Event   string `json:"event" bson:"event"`
	Payload string `json:"payload" bson:"payload"`
	ID      string `json:"id" bson:"id"`
}
