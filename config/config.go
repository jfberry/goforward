package config

type configDefinition struct {
	Port      int      `json:"port"`
	Webhooks  []string `json:"webhooks"`
	Controler string   `json:"controler"`
}

var Config configDefinition
