package lfi

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	lfi = "Local File Inclusion"
)

type LFIForm struct {
	common.PageQuery
	MethodRefer  string             `json:"mehod"`
	URL          string             `json:"url"`
	Param        string             `json:"param"`
	Cookie       string             `json:"cookie"`
	JWT          string             `json:"jwt"`
	ReportNumber primitive.ObjectID `json:"-"`
}

func (f LFIForm) URLOptions() *entities.LFI {
	webURL, err := url.Parse(f.URL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := utils.FetchParam(queries, f.Param)

	return &entities.LFI{
		URL:            webURL,
		Parameter:      p,
		ParameterValue: pValue,
	}
}
