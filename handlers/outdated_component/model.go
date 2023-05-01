package outdated_component

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"
)

const (
	lfi = "Local File Inclusion"
)

type OutdatedComponentForm struct {
	common.PageQuery
	MethodRefer string   `json:"mehod"`
	URL         string   `json:"url"`
	Param       string   `json:"param"`
	Cookie      string   `json:"cookie"`
	JWT         string   `json:"jwt"`
	Refer       []string `json:"-"`
}

func (f OutdatedComponentForm) URLOptions() *entities.OutdatedComponent {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := utils.FetchParam(queries, f.Param)

	return &entities.OutdatedComponent{
		URL:            webURL,
		Parameter:      p,
		ParameterValue: pValue,
	}
}
