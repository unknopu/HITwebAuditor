package entities

// PageInformation page information
type PageInformation struct {
	Vulnerabilities    int  `json:"total_number_of_vulnerability"`
	Low                int  `json:"total_number_of_low"`
	Medium             int  `json:"total_number_of_medium"`
	High               int  `json:"total_number_of_high"`
	Critical           int  `json:"total_number_of_critical"`
	RiskRate           int  `json:"risk_rate"`
	Injection          *int `json:"type_injection,omitempty"`
	Broken             *int `json:"type_broken_access_control,omitempty"`
	Cryptography       *int `json:"type_crypto_failure,omitempty"`
	MisConfiguration   *int `json:"type_miss_configuration,omitempty"`
	OutdatedComponents *int `json:"type_outdated_components,omitempty"`
}

// Page page model
type Page struct {
	PageInformation *PageInformation `json:"page_information,omitempty"`
	Entities        interface{}      `json:"entities,omitempty"`
}

// NewPage new page
func NewPage(pif PageInformation, es interface{}) *Page {
	return &Page{
		PageInformation: &PageInformation{
			Vulnerabilities: pif.Vulnerabilities,
			Low:             pif.Low,
			Medium:          pif.Medium,
			High:            pif.High,
			Critical:        pif.Critical,
			RiskRate:        calculateRiskRate(pif),
		},
		Entities: es,
	}
}

func calculateRiskRate(pif PageInformation) int {
	var numerator float64 = 0
	var denominator float64 = float64(pif.Vulnerabilities) * 4

	numerator += float64(pif.Low * 1)
	numerator += float64(pif.Medium * 2)
	numerator += float64(pif.High * 3)
	numerator += float64(pif.Critical * 4)

	return int((numerator / denominator) * 100)
}
