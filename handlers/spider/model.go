package spider

import "auditor/handlers/common"

type BaseForm struct {
	common.PageQuery
	BaseURL string `json:"base_url"`
}
