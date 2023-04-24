package utils

import "net/url"

func FetchParam(vs url.Values, param string) (string, string) {
	var key, value string
	for v := range vs {
		if param == v {
			return v, vs.Get(v)
		}
		key, value = v, vs.Get(v)
	}
	return key, value
}
