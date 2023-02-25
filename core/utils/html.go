package utils

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/k3a/html2text"
)

var (
	c = &tls.Config{InsecureSkipVerify: true}
)

func GetPageHTML(pageURL string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = c

	res, err := http.Get(pageURL)
	if err != nil || res.StatusCode == http.StatusForbidden {
		log.Fatal(err, res.StatusCode)
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	return string(html)
}

func GetPageLength(pageURL string) int {
	html := GetPageHTML(pageURL)
	if strings.Contains(html, "<head>") {
		afterHeadHTML := strings.SplitAfter(string(html), "<head>")
		plain := html2text.HTML2Text(afterHeadHTML[1])
		return len(plain)
	}
	return len(html)
}
