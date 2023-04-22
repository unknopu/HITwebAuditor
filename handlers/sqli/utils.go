package sqli

import (
	"auditor/core/utils"
	"auditor/entities"
)

func tryInjection(option entities.SQLi) bool {
	u := option.URL
	q := u.Query()
	q.Set(option.Parameter, option.ParameterValue+"'")
	u.RawQuery = q.Encode()

	len := utils.GetPageLength(u.String(), option.Cookie)

	return option.PageLength != len
}
