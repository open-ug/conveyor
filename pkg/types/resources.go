package types

type ResourceDefinition struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Schema      interface{} `json:"schema"`
}

type Resource struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Pipeline string            `json:"pipeline,omitempty"`
	Resource string            `json:"resource"`
	Metadata map[string]string `json:"metadata"`
	Spec     interface{}       `json:"spec"`
}
