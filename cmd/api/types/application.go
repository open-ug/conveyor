package types

type ApplicationSpec struct {
	AppName   string               `json:"app-name" bson:"app-name"`
	Image     string               `json:"image" bson:"image"`
	Volumes   []ApplicationVolume  `json:"volumes" bson:"volumes"`
	Ports     []ApplicationPortMap `json:"ports" bson:"ports"`
	Resources ApplicationResource  `json:"resources" bson:"resources"`
	Env       []ApplicationEnvVar  `json:"envFrom" bson:"envFrom"`
	Network   string               `json:"network" bson:"network"`
}
type ApplicationEnvVar struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type ApplicationResource struct {
	Storage int    `json:"storage" bson:"storage"`
	Memory  string `json:"memory" bson:"memory"`
	CPU     string `json:"cpu" bson:"cpu"`
}

type ApplicationVolume struct {
	VolumeName string `json:"volume-name" bson:"volume-name"`
	Path       string `json:"path" bson:"path"`
}

type ApplicationPortMap struct {
	Internal int    `json:"internal" bson:"internal"`
	External int    `json:"external" bson:"external"`
	Domain   string `json:"domain" bson:"domain"`
	SSL      bool   `json:"SSL" bson:"SSL"`
}

type Application struct {
	Name string          `json:"name" bson:"name"`
	Spec ApplicationSpec `json:"spec" bson:"spec"`
}

type ApplicationMsg struct {
	Action  string      `json:"action" bson:"action"`
	Payload Application `json:"payload" bson:"payload"`
	ID      string      `json:"id" bson:"id"`
}
