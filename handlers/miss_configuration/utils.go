package miss_configuration

import (
	"auditor/core/utils"
	"auditor/entities"
	"log"
	"net/http"
)

func buildPageInfomation(e []*entities.MissConfigurationReport) *entities.Page {
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

func fetchHeaders(option entities.MissConfiguration) *HttpHeader {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = utils.C

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, option.URL.String(), nil)

	res, err := client.Do(r)
	if err != nil {
		log.Println("[*] GET HTML: ", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	return &HttpHeader{
		Server:     res.Header.Get("Server"),
		XPoweredBy: res.Header.Get("X-Powered-By"),
	}
}

func anyVersionLeak(h *HttpHeader) bool {
	if h.Server != "" {
		return true
	}
	if h.XPoweredBy != "" {
		return true
	}
	return false
}
