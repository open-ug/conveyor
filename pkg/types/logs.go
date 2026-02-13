package types

type Log struct {
	// RunID is a unique identifier for the workflow run. This allows you to group logs by specific runs of your pipelines.
	RunID string `json:"runid"`
	// Driver is the name of the driver that generated the log. This helps in identifying which component of the system produced the log.
	Driver string `json:"driver"`
	// Pipeline is the name of the pipeline associated with the log. This allows you to trace logs back to specific pipelines.
	Pipeline string `json:"pipeline"`
	// Timestamp is the time when the log was created. This is useful for ordering logs and understanding the sequence of events.
	Timestamp string `json:"timestamp"`
	// Message is the actual log message. This contains the information or error details that you want to record.
	Message string `json:"message"`
}
