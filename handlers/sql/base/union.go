package base

import (
	"auditor/core/utils"
	"auditor/entities"
	"fmt"
	"log"
	"strings"
)

const (
	ErrXPathForm = "XPATH syntax error: '.*'"
	ErrXPathQueryFrom     = "XPATH syntax error: ':.*'"
)

func UnionBasedvalidate(options *entities.DBOptions, query string) string {
	u := *options.URL
	u.RawQuery = fmt.Sprintf("%s=%s", options.Parameter, options.ParameterValue+query)

	log.Println("========================")
	log.Println(options.Parameter)
	log.Println(options.ParameterValue + query)
	log.Println(u.RawQuery)
	log.Println(u.String())
	log.Println("========================")

	return utils.GetPageHTML(u.String(), options.Cookie)
}

func TrimData(s string) string {
	s = strings.ReplaceAll(s, "XPATH syntax error: ':", "")
	s = strings.ReplaceAll(s, "'\nWarning:", "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}
