package utils

import (
	"auditor/core/mongodb"
	"log"
	"net/http"
	"net/url"
)

type BasicRequestForm struct {
	mongodb.Model  `bson:",inline"`
	URL            *url.URL `json:"url,omitempty"`
	Parameter      string   `json:"param,omitempty"`
	ParameterValue string   `json:"param_value,omitempty"`
}

func SendRequest(option BasicRequestForm, payload string) *http.Response {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = C

	if payload != "" {
		q := option.URL.Query()
		q.Set(option.Parameter, payload)
		option.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, option.URL.String(), nil)

	res, err := client.Do(r)
	if err != nil {
		log.Println("[*] GET HTML: ", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	return res
}
