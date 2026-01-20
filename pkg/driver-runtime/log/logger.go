package log

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/pkg/types"
)

type DriverLogger struct {
	// The driver logger is responsible for logging the driver
	// and the driver lifecycle.
	// Logger is the logger for the driver manager
	Logger Logger

	// The driver name
	DriverName string

	// Labels are the labels that will be used to identify the logs
	Labels map[string]string

	// NatsCon is the nats client for the driver logger
	NatsCon *nats.Conn
}

type Logger interface {
	PushLog(logEntry types.Log) error
}

func NewDriverLogger(driverName string, labels map[string]string, natsCon *nats.Conn) *DriverLogger {

	return &DriverLogger{
		DriverName: driverName,
		Logger:     NewDefaultLogger(natsCon),
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

	logEntry := types.Log{
		RunID:     initialLabels["run_id"],
		Driver:    initialLabels["driver"],
		Message:   message,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	// Send the log to Loki
	err := d.Logger.PushLog(logEntry)
	if err != nil {
		return err
	}
	// Marshal the message to JSON
	messageBytes, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	// Send the log to Nats
	err = d.NatsCon.Publish("live.logs."+runId+"."+d.DriverName, messageBytes)

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

type DefaultLogger struct {
	NatsCon *nats.Conn
	js      jetstream.JetStream
}

func NewDefaultLogger(natsCon *nats.Conn) *DefaultLogger {
	js, err := jetstream.New(natsCon)
	if err != nil {
		color.Red("Error Occured while creating JetStream context: %v", err)
		return nil
	}
	return &DefaultLogger{
		NatsCon: natsCon,
		js:      js,
	}
}

func (d *DefaultLogger) PushLog(logEntry types.Log) error {

	// Marshal the log entry to JSON
	logEntryBytes, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	// Publish the log entry to NATS JetStream
	_, err = d.js.Publish(context.Background(), "logs."+logEntry.RunID, logEntryBytes)
	if err != nil {
		return err
	}

	return nil
}
