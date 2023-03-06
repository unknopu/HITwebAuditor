package sql

import "auditor/handlers/common"

type BaseForm struct {
	common.PageQuery
	URL    string  `json:"url"`
	Param  string  `json:"param"`
	Level  *string `json:"level"`
	Cookie string  `json:"cookie"`
}
