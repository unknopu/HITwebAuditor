package entities

import (
	"auditor/core/mongodb"
	"auditor/core/utils"
	"net/url"
)

type DBOptions struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL
	Parameter      string
	ParameterValue string
	PageLength     int    `json:"-"`
	PageOrigin     string `json:"-"`
	Payload        int    `json:"-"`
	PayloadStr     string `json:"-"`
	NameLength     int    `json:"-"`
	TableCount     int
	Tables         map[string][]string
	Name           string
	Rows           map[int][]string
	Cookie         string `json:"cookie,omitempty"`
	FromDB         bool   `bson:"-"`
	InjectionBase  string
}

func URLOptions(rURL, param, cookie string) *DBOptions {
	webURL, err := url.Parse(rURL)
	if err != nil {
		return nil
	}

	queries, _ := url.ParseQuery(webURL.RawQuery)
	p, pValue := fetchParam(queries, param)

	return &DBOptions{
		URL:            webURL,
		PageOrigin:     utils.GetPageHTML(webURL.String(), cookie),
		PageLength:     utils.GetPageLength(webURL.String(), cookie),
		Tables:         make(map[string][]string),
		Rows:           make(map[int][]string),
		Parameter:      p,
		ParameterValue: pValue,
		Cookie:         cookie,
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

type ProcType int

const (
	NameLength ProcType = iota
	Name
	TableCount
	ColumnsName
	Tables
)

func (i DBOptions) ValidateProc(proc ProcType) bool {
	if proc == NameLength {
		return i.NameLength > 0
	}
	if proc == Name {
		return len(i.Name) == i.NameLength
	}
	if proc == TableCount {
		return i.TableCount > 0
	}
	if proc == Tables || proc == ColumnsName {
		return len(i.Tables) == i.TableCount
	}

	return false
}
