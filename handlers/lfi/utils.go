package lfi

import (
	"auditor/core/utils"
	"auditor/entities"
	"io/ioutil"
	"log"
	"net/http"
)

func buildPageInfomation(e []*entities.LFIReport) *entities.Page {
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

func injectPayload(option entities.LFI, payload string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = utils.C

	q := option.URL.Query()
	q.Set(option.Parameter, payload)
	option.URL.RawQuery = q.Encode()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, option.URL.String(), nil)

	res, err := client.Do(r)
	if err != nil {
		log.Println("[*] GET HTML: ", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
