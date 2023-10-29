package models

type StandInfo struct {
	Components      []ComponentInfo `json:"components,omitempty"`
	TotalComponents int             `json:"total_components,omitempty"`
}

type ComponentInfo struct {
	Name       string      `json:"name,omitempty"`
	Containers []Container `json:"containers,omitempty"`
	Address    string      `json:"address,omitempty"`
	Status     string      `json:"status,omitempty"`
}

type Container struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Ready bool   `json:"ready,omitempty"`
}
