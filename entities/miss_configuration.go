package entities

import (
	"auditor/core/mongodb"
	"net/url"
)

type MissConfiguration struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL `json:"url,omitempty"`
	Parameter      string   `json:"param,omitempty"`
	ParameterValue string   `json:"param_value,omitempty"`
}

type MissConfigurationReport struct {
	Location       string          `json:"location,omitempty"`
	Payload        []string        `json:"payload,omitempty"`
	Level          []string        `json:"level,omitempty"`
	Type           TYPE            `json:"type,omitempty"`
	Vaulnerability []VULNERABILITY `json:"vaulnerability,omitempty"`
}
