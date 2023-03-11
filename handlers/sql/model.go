package sql

import "auditor/handlers/common"

type BaseForm struct {
	common.PageQuery
	URL        string  `json:"url"`
	Param      string  `json:"param"`
	Level      *string `json:"level"`
	QueryTable string  `json:"query_table"`
	Cookie     string  `json:"cookie"`
}
