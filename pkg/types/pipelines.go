package types

type Pipeline struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Resource    string            `json:"resource"`
	Steps       []Step            `json:"steps"`
	Metadata    map[string]string `json:"metadata"`
}

type Step struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

type DriverResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
