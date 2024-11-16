package utils

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	C = &tls.Config{InsecureSkipVerify: true}
)

func GetPageHTML(pageURL, cookie string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = C

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, pageURL, nil)
	r.Header.Add("Cookie", cookie)

	res, err := client.Do(r)
	if err != nil {
		log.Println("[*] GET HTML: ", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(html)
}

func GetPageLength(pageURL, cookie string) int {
	html := GetPageHTML(pageURL, cookie)
	return len(html)
}

func ISContainsHTMLBody(htmlString string) bool {
	re := regexp.MustCompile(`(?i)<body[^>]*>.*<\/body>`)
	return re.MatchString(htmlString)
}
