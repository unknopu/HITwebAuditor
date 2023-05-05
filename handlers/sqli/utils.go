package sqli

import (
	"auditor/core/utils"
	"auditor/entities"
	"fmt"
	"regexp"
	"strings"
)

func buildPageInfomation(e []*entities.SQLiReport) *entities.Page {
	if len(e) <= 0 {
		return nil
	}

	var low, medium, high, critical int

	for _, r := range e {
		switch r.Level {
		case entities.LOW:
			low += 1
		case entities.MEDIUM:
			medium += 1
		case entities.HIGH:
			high += 1
		case entities.CRITICAL:
			critical += 1
		}
	}

	pif := &entities.PageInformation{
		Vulnerabilities: len(e),
		Low:             low,
		Medium:          medium,
		High:            high,
		Critical:        critical,
	}

	return entities.NewPage(*pif, e)
}

func isParamInjectable(option entities.SQLi) (bool, string) {
	u := option.URL
	q := u.Query()
	q.Set(option.Parameter, option.ParameterValue+"'")
	u.RawQuery = q.Encode()

	body := utils.GetPageHTML(u.String(), option.Cookie)

	return option.PageLength != len(body), body
}

func detectErrMsg(httpBody string) bool {
	for _, payload := range errPayloads {
		if strings.Contains(httpBody, payload) {
			return true
		}
	}
	return false
}

func isContainBooleanBased(option entities.SQLi) bool {
	for _, payload := range booleanPayloads {
		u := option.URL
		q := u.Query()
		q.Set(option.Parameter, option.ParameterValue+payload)
		u.RawQuery = q.Encode()
		body := utils.GetPageHTML(u.String(), option.Cookie)

		if !detectErrMsg(body) {
			return strings.Contains(body, `<body>`)
		}
	}

	return false
}

func isContainUnionBased(option entities.SQLi) bool {
	u := option.URL
	var column int
	for column = 1; column < 50; column++ {

		payload := fmt.Sprintf("+order+by+%v", column)
		u.RawQuery = fmt.Sprintf("%v=%v", option.Parameter, option.ParameterValue+payload)
		body := utils.GetPageHTML(u.String(), option.Cookie)

		if detectErrMsg(body) {
			payload = buildUnionPayload(column - 1)
			u.RawQuery = fmt.Sprintf("%v=%v", option.Parameter, payload)

			body := utils.GetPageHTML(u.String(), option.Cookie)
			matched, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, body)

			return matched
		}
	}

	return false
}

func buildUnionPayload(n int) string {
	payload := "9999+union+select+1"
	for i := 2; i <= n; i++ {
		payload += ",NOW()"
	}
	return payload
}
