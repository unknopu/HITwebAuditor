package report

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Form struct {
	common.PageQuery
	MethodRefer  string             `json:"mehod"`
	URL          string             `json:"url"`
	Param        string             `json:"param"`
	Cookie       string             `json:"cookie"`
	JWT          string             `json:"jwt"`
	ReportNumber primitive.ObjectID `json:"-"`
}

func (f Form) URLOptions() *entities.ReportBase {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := utils.FetchParam(queries, f.Param)

	return &entities.ReportBase{
		URL:            webURL,
		Parameter:      p,
		ParameterValue: pValue,
	}
}

type GetLatestForm struct {
	common.PageQuery
}
