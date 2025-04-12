package log

import (
	lokiLog "github.com/open-ug/conveyor/internal/logging"
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
}

func NewDriverLogger(driverName string, labels map[string]string) *DriverLogger {
	// Load the configuration
	lokiClient := lokiLog.New("http://localhost:3100")

	return &DriverLogger{
		DriverName: driverName,
		Logger:     lokiClient,
		Labels:     labels,
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
	// Send the log to Loki
	err := d.Logger.PushLog(initialLabels, message)
	if err != nil {
		return err
	}

	return nil
}
