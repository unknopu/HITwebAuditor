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

func GetPageHTML(pageURL, cookie string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = c

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
	if strings.Contains(html, "<head>") {
		afterHeadHTML := strings.SplitAfter(string(html), "<head>")
		plain := html2text.HTML2Text(afterHeadHTML[1])
		return len(plain)
	}
	return len(html)
}

func setCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:  "Cookie",
		Value: value,
	}
}
