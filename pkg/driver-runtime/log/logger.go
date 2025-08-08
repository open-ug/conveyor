package log

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/spf13/viper"
)

type DriverLogger struct {
	// The driver logger is responsible for logging the driver
	// and the driver lifecycle.
	// Logger is the logger for the driver manager
	Logger *utils.LokiClient

	// The driver name
	DriverName string

	// Labels are the labels that will be used to identify the logs
	Labels map[string]string

	// NatsCon is the nats client for the driver logger
	NatsCon *nats.Conn
}

func NewDriverLogger(driverName string, labels map[string]string, natsCon *nats.Conn) *DriverLogger {
	// Load the configuration
	lokiClient := utils.NewLokiClient(
		viper.GetString("loki.host"),
	)

	return &DriverLogger{
		DriverName: driverName,
		Logger:     lokiClient,
		Labels:     labels,
		NatsCon:    natsCon,
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

	// Send the log to Nats
	err = d.NatsCon.Publish("driver:"+d.DriverName+":logs:"+runId, messageBytes)

	if err != nil {
		return err
	}

	return nil
}

func (d *DriverLogger) Write(p []byte) (n int, err error) {
	// Send the log to Loki
	err = d.Log(nil, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
