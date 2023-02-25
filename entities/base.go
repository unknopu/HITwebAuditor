package entities

import (
	"auditor/core/utils"
	"net/url"
)

type DBOptions struct {
	URL            *url.URL
	Parameter      string
	ParameterValue string
	PageLength     int
	Payload        int
	NameLength     int
	TableCount     int
	Name           string
	TablesColumns  map[string][]string
	TablesRows     map[int][]string
}

func URLOptions(rURL, param string) *DBOptions {
	webURL, err := url.Parse(rURL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := fetchParam(queries, param)

	return &DBOptions{
		URL:            webURL,
		PageLength:     utils.GetPageLength(webURL.String()),
		TablesColumns:  make(map[string][]string),
		TablesRows:     make(map[int][]string),
		Parameter:      p,
		ParameterValue: pValue,
	}
}

func fetchParam(vs url.Values, param string) (string, string) {
	var key, value string
	for v := range vs {
		if param == v {
			return v, vs.Get(v)
		}
		key, value = v, vs.Get(v)
	}
	return key, value
}
