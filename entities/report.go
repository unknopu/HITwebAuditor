package entities

import (
	"auditor/core/mongodb"
	"encoding/json"
	"net/url"
)

type ReportBase struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL `json:"url,omitempty"`
	Parameter      string   `json:"param,omitempty"`
	ParameterValue string   `json:"param_value,omitempty"`
}

type Report struct {
	URL               string `json:"url"`
	RiskRate          int    `json:"risk_rate"`
	SQLi              *Page  `json:"sqli,omitempty"`
	LFI               *Page  `json:"lfi,omitempty"`
	MConfig           *Page  `json:"miss_config,omitempty"`
	XSS               *Page  `json:"xss,omitempty"`
	CryptoFailure     *Page  `json:"crypto_failure,omitempty"`
	OutdatedComponent *Page  `json:"outdated_component,omitempty"`
}

// MarshalJSON custom image json
func (i Report) MarshalJSON() ([]byte, error) {
	type Alias Report
	m := &struct {
		PageInformation PageInformation `json:"page_information,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&i),
	}

	var counter, raiskRate, vulnerabilities, low, medium, high, critical int
	if m.SQLi != nil {
		counter++
		raiskRate += m.SQLi.PageInformation.RiskRate
		vulnerabilities += m.SQLi.PageInformation.Vulnerabilities
		low += m.SQLi.PageInformation.Low
		medium += m.SQLi.PageInformation.Medium
		high += m.SQLi.PageInformation.High
		critical += m.SQLi.PageInformation.Critical
	}

	if m.LFI != nil {
		counter++
		raiskRate += m.LFI.PageInformation.RiskRate
		vulnerabilities += m.LFI.PageInformation.Vulnerabilities
		low += m.LFI.PageInformation.Low
		medium += m.LFI.PageInformation.Medium
		high += m.LFI.PageInformation.High
		critical += m.LFI.PageInformation.Critical
	}

	if m.MConfig != nil {
		counter++
		raiskRate += m.MConfig.PageInformation.RiskRate
		vulnerabilities += m.MConfig.PageInformation.Vulnerabilities
		low += m.MConfig.PageInformation.Low
		medium += m.MConfig.PageInformation.Medium
		high += m.MConfig.PageInformation.High
		critical += m.MConfig.PageInformation.Critical
	}

	if m.XSS != nil {
		counter++
		raiskRate += m.XSS.PageInformation.RiskRate
		vulnerabilities += m.XSS.PageInformation.Vulnerabilities
		low += m.XSS.PageInformation.Low
		medium += m.XSS.PageInformation.Medium
		high += m.XSS.PageInformation.High
		critical += m.XSS.PageInformation.Critical
	}

	if m.CryptoFailure != nil {
		counter++
		raiskRate += m.CryptoFailure.PageInformation.RiskRate
		vulnerabilities += m.CryptoFailure.PageInformation.Vulnerabilities
		low += m.CryptoFailure.PageInformation.Low
		medium += m.CryptoFailure.PageInformation.Medium
		high += m.CryptoFailure.PageInformation.High
		critical += m.CryptoFailure.PageInformation.Critical
	}

	if m.OutdatedComponent != nil {
		counter++
		raiskRate += m.OutdatedComponent.PageInformation.RiskRate
		vulnerabilities += m.OutdatedComponent.PageInformation.Vulnerabilities
		low += m.OutdatedComponent.PageInformation.Low
		medium += m.OutdatedComponent.PageInformation.Medium
		high += m.OutdatedComponent.PageInformation.High
		critical += m.OutdatedComponent.PageInformation.Critical
	}

	m.PageInformation.Vulnerabilities = vulnerabilities
	m.PageInformation.Low = low
	m.PageInformation.Medium = medium
	m.PageInformation.High = high
	m.PageInformation.Critical = critical
	m.PageInformation.RiskRate = int((raiskRate / counter))

	return json.Marshal(m)
}
