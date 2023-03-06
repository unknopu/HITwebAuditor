package base

import (
	"auditor/core/utils"
	"auditor/entities"
	"fmt"
	"log"
)

func UnionBasedvalidate(options *entities.DBOptions, query string) string {
	u := *options.URL
	// q := u.Query()
	// q.Set(options.Parameter, options.ParameterValue+query)
	u.RawQuery = fmt.Sprintf("%s=%s", options.Parameter, options.ParameterValue+query)

	log.Println("========================")
	log.Println(options.Parameter)
	log.Println(options.ParameterValue + query)
	log.Println(u.RawQuery)
	log.Println(u.String())
	log.Println("========================")

	return utils.GetPageHTML(u.String(), options.Cookie)
}
