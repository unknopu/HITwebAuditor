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
	PageLength     int
	Payload        int
	NameLength     int
	TableCount     int
	Name           string
	Columns        map[string][]string
	Rows           map[int][]string
	Cookie         string
	FromDB         bool `bson:"-"`
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
		PageLength:     utils.GetPageLength(webURL.String(), cookie),
		Columns:        make(map[string][]string),
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
		return len(i.Name) > 0
	}
	if proc == TableCount {
		return i.TableCount > 0
	}
	if proc == ColumnsName {
		return len(i.Columns) > 0
	}
	if proc == Tables {
		return i.Columns == nil
	}
	return false
}
