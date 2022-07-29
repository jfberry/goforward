package config

type configDefinition struct {
	Port     int      `json:"port"`
	Webhooks []string `json:"webhooks"`
}

var Config configDefinition
