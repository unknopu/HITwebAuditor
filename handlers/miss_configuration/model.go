package miss_configuration

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"
)

const (
	SERVER_RULE   = `^(nginx|Nginx|apache|Apache)/(\d+\.)?(\d+\.)?(\*|\d+)`
	POWERDBY_RULE = `^(php|PHP)/(\d+\.)?(\d+\.)?(\*|\d+)`
)

type MCForm struct {
	common.PageQuery
	MethodRefer string `json:"mehod"`
	URL         string `json:"url"`
	Param       string `json:"param"`
	Cookie      string `json:"cookie"`
	JWT         string `json:"jwt"`
}

func (f MCForm) URLOptions() *entities.MissConfiguration {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := utils.FetchParam(queries, f.Param)

	return &entities.MissConfiguration{
		URL:            webURL,
		Parameter:      p,
		ParameterValue: pValue,
	}
}

type HttpHeader struct {
	Server     string `json:"server"`
	XPoweredBy string `json:"x_powered_by"`
}
