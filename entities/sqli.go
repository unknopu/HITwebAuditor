package entities

import (
	"auditor/core/mongodb"
	"net/url"
)

type SQLi struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL `json:"url,omitempty"`
	Parameter      string   `json:"param,omitempty"`
	ParameterValue string   `json:"param_value,omitempty"`
	PageLength     int      `json:"page_length,omitempty"`
	PageOrigin     string   `json:"-"`
	Payload        int      `json:"payload,omitempty"`
	PayloadStr     string   `json:"payload_str,omitempty"`
	NameLength     int      `json:"name_length,omitempty"`
	TableCount     int      `json:"table_count,omitempty"`
	Tables         []string `json:"table,omitempty"`
	Rows           []string `json:"rows,omitempty"`
	Cookie         string   `json:"cookie,omitempty"`
	DatabaseName   bool     `bson:"database_name,omitempty"`
	Base           []string `json:"base,omitempty"`
}

type SQLiReport struct {
	mongodb.Model  `bson:",inline"`
	Location       string        `json:"location,omitempty"`
	Payload        []string      `json:"payload,omitempty"`
	Level          LEVEL         `json:"level,omitempty"`
	Type           TYPE          `json:"type,omitempty"`
	Vaulnerability VULNERABILITY `json:"vaulnerability,omitempty"`
}
