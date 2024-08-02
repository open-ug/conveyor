package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Application `json:"items"`
}

type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ApplicationSpec `json:"spec"`
}

type ApplicationSpec struct {
	AppName   string               `json:"app-name"`
	Image     string               `json:"image"`
	Volumes   []ApplicationVolume  `json:"volumes"`
	Ports     []ApplicationPortMap `json:"ports"`
	EnvFrom   string               `json:"envFrom"`
	Resources ApplicationResource  `json:"resources"`
}

type ApplicationResource struct {
	Storage int    `json:"storage"`
	Memory  string `json:"memory"`
	CPU     string `json:"cpu"`
}

type ApplicationVolume struct {
	VolumeName string `json:"volume-name"`
	Path       string `json:"path"`
}

type ApplicationPortMap struct {
	Internal int    `json:"internal"`
	External int    `json:"external"`
	Domain   string `json:"domain"`
	SSL      bool   `json:"SSL"`
}
