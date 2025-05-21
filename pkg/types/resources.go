package types

type ResourceDefinition struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Schema      interface{} `json:"schema"`
}

type Resource struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Resource string      `json:"resource"`
	Spec     interface{} `json:"spec"`
}
