package sql

import (
	"auditor/core/utils"
	based "auditor/handlers/sql/base"
	"auditor/payloads/intruder/detect"

	"strings"
)

func validateByErrorBased() int {
	u := *options.URL
	q := u.Query()
	q.Set(options.Parameter, options.ParameterValue+"'")
	u.RawQuery = q.Encode()
	html := utils.GetPageHTML(u.String(), options.Cookie)
	for _, valueErr := range detect.ErrPayloads {
		if !strings.Contains(html, valueErr) {
			return 1
		}
	}
	return 0
}

func validateByMethod(query string, method based.SQLi) int {
	u := *options.URL
	q := u.Query()
	q.Set(options.Parameter, options.ParameterValue+query)
	u.RawQuery = q.Encode()

	switch method {
	case based.LengthValidation:
		secondLen := utils.GetPageLength(u.String(), options.Cookie)
		if options.PageLength == secondLen {
			return 1
		}
		return 0

	case based.ErrorSQLiBased:
		html := utils.GetPageHTML(u.String(), options.Cookie)
		for _, valueErr := range detect.ErrPayloads {
			if !strings.Contains(html, valueErr) {
				return 1
			}
		}
		return 0
	}

	return 0
}
