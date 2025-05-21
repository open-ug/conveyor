/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package types

// Application Object
type Application struct {
	Name string          `json:"name" bson:"name"`
	Spec ApplicationSpec `json:"spec" bson:"spec"`
}

// Application Object Spec
type ApplicationSpec struct {
	AppName   string               `json:"app-name" bson:"app-name"`
	Source    ApplicationSource    `json:"source" bson:"source"`
	Volumes   []ApplicationVolume  `json:"volumes" bson:"volumes"`
	Ports     []ApplicationPortMap `json:"ports" bson:"ports"`
	Resources ApplicationResource  `json:"resources" bson:"resources"`
	Env       []ApplicationEnvVar  `json:"envFrom" bson:"envFrom"`
	Network   string               `json:"network" bson:"network"`
}

// Application Environment Variable
type ApplicationEnvVar struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Application Resource Object. It defines the CPU and Memory requirements for the application
type ApplicationResource struct {
	// Storage is the amount of storage to be allocated to the application
	Storage int `json:"storage" bson:"storage"`
	// Memory is the amount of memory to be allocated to the application
	Memory string `json:"memory" bson:"memory"`
	// CPU is the amount of CPU to be allocated to the application
	CPU string `json:"cpu" bson:"cpu"`
}

// Application Volume Object. It defines the volume name and path to be mounted to the application
type ApplicationVolume struct {
	// Name of the volume to be mounted to the application
	VolumeName string `json:"volume-name" bson:"volume-name"`
	// Path is the path in the container where the volume will be mounted
	Path string `json:"path" bson:"path"`
}

// The Port Map Object. It defines the internal and external ports to be mapped to the application and if Domains or SSL are to be used
type ApplicationPortMap struct {
	// Internal is the port exposed by the application in the container
	Internal int `json:"internal" bson:"internal"`
	// External is the port exposed by the application on the host
	External int `json:"external" bson:"external"`
	// Domain is the domain name to be used for the traffic to this Specific Port.
	Domain string `json:"domain" bson:"domain"`
	// SSL is a boolean value that determines if SSL is to be used for the traffic to this Specific Port.
	SSL bool `json:"SSL" bson:"SSL"`
}

type ApplicationMsg struct {
	Action  string      `json:"action" bson:"action"`
	Payload Application `json:"payload" bson:"payload"`
	ID      string      `json:"id" bson:"id"`
}

// Source Object for Git Repositories
type GitRepo struct {
	// URL is the URL of the Git Repository
	URL string `json:"url" bson:"url"`
	// Branch is the branch of the Git Repository
	Branch string `json:"branch" bson:"branch"`
	// Revision is the revision of the Git Repository
	Revision string `json:"revision" bson:"revision"`
	// Username is the username for the Git Repository
	Username string `json:"username" bson:"username"`
	// Password is the password for the Git Repository
	Password string `json:"password" bson:"password"`
}

// Source Object for Blob Files
type BlobFileSource struct {
	Source string `json:"source" bson:"source"`
}

// Source Object for Container Images
type ImageSource struct {
	ImageURI string `json:"imageURI" bson:"imageURI"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Source Object for Applications
type ApplicationSource struct {
	Type     string         `json:"type" bson:"type"`
	GitRepo  GitRepo        `json:"gitRepo" bson:"gitRepo"`
	BlobFile BlobFileSource `json:"blobFile" bson:"blobFile"`
	Image    ImageSource    `json:"image" bson:"image"`
}
