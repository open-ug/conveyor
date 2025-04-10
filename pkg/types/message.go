package types

type DriverMessage struct {
	// Event Name rg. `docker-build-complete`
	Event string `json:"event" bson:"event"`
	// JSON Payload
	Payload string `json:"payload" bson:"payload"`
	ID      string `json:"id" bson:"id"`
}
