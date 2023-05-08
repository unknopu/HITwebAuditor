package static_analysis

import (
	"auditor/handlers/common"
	"mime/multipart"
)

type StaticAnalysisForm struct {
	common.PageQuery
	URL  *string         `json:"url" form:"url"`
	File *multipart.File `json:"file" form:"file"`
}

// func (f SqliForm) URLOptions() *entities.SQLi {
// 	webURL, err := url.Parse(f.URL)
// 	if err != nil {
// 		return nil
// 	}

// 	queries, _ := url.ParseQuery(webURL.RawQuery)
// 	p, pValue := utils.FetchParam(queries, f.Param)

// 	// log.Println("url origin: ", webURL, p, pValue)
// 	return &entities.SQLi{
// 		URL:            webURL,
// 		PageOrigin:     utils.GetPageHTML(webURL.String(), f.Cookie),
// 		PageLength:     utils.GetPageLength(webURL.String(), f.Cookie),
// 		Parameter:      p,
// 		ParameterValue: pValue,
// 		Cookie:         f.Cookie,
// 	}
// }
