package sqli

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SqliForm struct {
	common.PageQuery
	MethodRefer  string             `json:"mehod"`
	URL          string             `json:"url"`
	Param        string             `json:"param"`
	Cookie       string             `json:"cookie"`
	JWT          string             `json:"jwt"`
	ReportNumber primitive.ObjectID `json:"-"`
}

func (f SqliForm) URLOptions() *entities.SQLi {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := utils.FetchParam(queries, f.Param)

	// log.Println("url origin: ", webURL, p, pValue)
	return &entities.SQLi{
		URL:            webURL,
		PageOrigin:     utils.GetPageHTML(webURL.String(), f.Cookie),
		PageLength:     utils.GetPageLength(webURL.String(), f.Cookie),
		Parameter:      p,
		ParameterValue: pValue,
		Cookie:         f.Cookie,
	}
}
