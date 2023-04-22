package lfi


import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"
)

const (
	Boolean = "Boolean Based SQL Injection"
)

type SqliForm struct {
	common.PageQuery
	MethodRefer string `json:"mehod"`
	URL         string `json:"url"`
	Param       string `json:"param"`
	Cookie      string `json:"cookie"`
	JWT         string `json:"jwt"`
}

func (f SqliForm) URLOptions() *entities.SQLi {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := fetchParam(queries, f.Param)

	return &entities.SQLi{
		URL:            webURL,
		PageOrigin:     utils.GetPageHTML(webURL.String(), f.Cookie),
		PageLength:     utils.GetPageLength(webURL.String(), f.Cookie),
		Parameter:      p,
		ParameterValue: pValue,
		Cookie:         f.Cookie,
	}
}

func fetchParam(vs url.Values, param string) (string, string) {
	var key, value string
	for v := range vs {
		if param == v {
			return v, vs.Get(v)
		}
		key, value = v, vs.Get(v)
	}
	return key, value
}
