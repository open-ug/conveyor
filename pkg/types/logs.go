package types

type Log struct {
	RunID     string `json:"runid"`
	Driver    string `json:"driver"`
	Pipeline  string `json:"pipeline"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
