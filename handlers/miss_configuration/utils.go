package miss_configuration

import (
	"auditor/core/utils"
	"auditor/entities"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func checkPWNVersion(version string, minMajor, minMinor, maxMajor, maxMinor int) bool {
	// Split the version string into major, minor, and patch components
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
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
	if major <= minMajor && (minor < minMinor) { // major == 5 && minor < 3
		return false // Version is less than 5.3.0
	}
	if major >= maxMajor && (minor > maxMinor) { //major == 7 && minor > 0
		return false // Version is greater than 7.0.0
	}
	log.Println("test: ", version)
	return true // Version is within the range 5.3.0 to 7.0.0
}
