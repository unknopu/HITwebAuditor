package entities

import (
	"auditor/core/mongodb"
	"net/url"
)

type XSS struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL `json:"url,omitempty"`
	Parameter      string   `json:"param,omitempty"`
	ParameterValue string   `json:"param_value,omitempty"`
}

type XSSReport struct {
	mongodb.Model  `bson:",inline"`
	Location       string        `json:"location,omitempty"`
	Payload        []string      `json:"payload,omitempty"`
	Level          LEVEL         `json:"level,omitempty"`
	Type           TYPE          `json:"type,omitempty"`
	Vaulnerability VULNERABILITY `json:"vaulnerability,omitempty"`
}
