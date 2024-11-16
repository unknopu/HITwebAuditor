package outdated_component

import (
	"auditor/core/utils"
	"auditor/entities"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func BuildPageInfomation(e []*entities.OutdatedComponentsReport) *entities.Page {
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

func isPhpPWN(version string) bool {
	// Split the version string into major, minor, and patch components
	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		return false // Version string is not in the correct format
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false // Major component is not a valid integer
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false // Minor component is not a valid integer
	}
	sub, err := strconv.Atoi(parts[2])
	if err != nil {
		return false // Minor component is not a valid integer
	}
	if major <= 5 && (minor < 3) {
		return false // Version is less than 5.3.0
	}
	if major >= 7 && (minor >= 0) && sub > 0 {
		return false // Version is greater than 7.0.0
	}
	return true // Version is within the range 5.3.0 to 7.0.0
}

func fetchHeaders(option entities.OutdatedComponent) *HttpHeader {
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
