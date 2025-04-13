package log

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	lokiLog "github.com/open-ug/conveyor/internal/logging"
	"github.com/redis/go-redis/v9"
)

type DriverLogger struct {
	// The driver logger is responsible for logging the driver
	// and the driver lifecycle.
	// Logger is the logger for the driver manager
	Logger *lokiLog.LokiClient

	// The driver name
	DriverName string

	// Labels are the labels that will be used to identify the logs
	Labels map[string]string

	// RedisClient is the redis client for the driver logger
	RedisClient *redis.Client
}

func NewDriverLogger(driverName string, labels map[string]string, redisClient *redis.Client) *DriverLogger {
	// Load the configuration
	lokiClient := lokiLog.New("http://localhost:3100")

	return &DriverLogger{
		DriverName:  driverName,
		Logger:      lokiClient,
		Labels:      labels,
		RedisClient: redisClient,
	}
}

func (d *DriverLogger) Log(labels map[string]string, message string) error {
	// Send the log to Loki
	initialLabels := map[string]string{
		"driver": d.DriverName,
	}
	// Add the labels to the initial labels
	for k, v := range d.Labels {
		initialLabels[k] = v
	}
	// Merge the labels with the initial labels
	for k, v := range labels {
		initialLabels[k] = v
	}
	runId := d.Labels["run_id"]
	// Send the log to Loki
	err := d.Logger.PushLog(initialLabels, message)
	if err != nil {
		return err
	}
	// an array of the current timestamp and message
	timestamp := []string{strconv.FormatInt(time.Now().Unix(), 10), message}
	// Marshal the message to JSON
	messageBytes, err := json.Marshal(timestamp)
	if err != nil {
		return err
	}
	// Convert the message to a string
	msg := string(messageBytes)
	// Send the log to Redis
	err = d.RedisClient.Publish(context.Background(), "driver:"+d.DriverName+":logs:"+runId, msg).Err()

	return nil
}
