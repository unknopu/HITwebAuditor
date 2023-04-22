package spider

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func spider(urlStr string, visited map[string]bool, depth int, swap *[]string) {
	if depth == 0 {
		return
	}

	visited[urlStr] = true

	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	links := getLinks(resp, urlStr)
	for _, link := range links {
		*swap = append(*swap, link)
		if !visited[link] {
			spider(link, visited, depth-1, swap)
		}
	}
}

func getLinks(body *http.Response, base string) []string {
	links := []string{}

	baseURL, _ := url.Parse(base)

	tokenizer := html.NewTokenizer(body.Body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					linkURL, err := baseURL.Parse(attr.Val)
					if err == nil && linkURL.Host == baseURL.Host {
						links = append(links, strings.TrimSuffix(linkURL.String(), "/"))
					}
				}
			}
		}
	}
}

func removeDuplicateLink(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}


// docker run --tty --interactive -d -p 8080:80 -v /Users/lin/Documents/2023大四春季学期/resecrhing/nginx:/usr/share/nginx/html nginx
// docker run --tty --interactive -d -p 8080:80 nginx