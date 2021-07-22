package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func getUrlResponse(u string) (*http.Response, error) {
	return http.Get(u)
}

func getUrlsFromBody(resp *http.Response, respErr error) []string {

	if respErr != nil {
		log.Fatalf("something went wrong when running: %s \n", respErr)
	}

	defer resp.Body.Close()
	doc, parseErr := goquery.NewDocumentFromReader(resp.Body)

	if parseErr != nil {
		log.Fatalf("something went wrong while parsing: %s \n", parseErr)
	}

	urls := []string{}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if hasHttpPrefix(url) {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		log.Fatalf("Could not find any links: %s \n", resp.Request.URL)
	}

	return urls
}

func hasHttpPrefix(url string) bool {
	return strings.HasPrefix(url, "http")
}

func getInfoAndSelelectedHeaders(resp *http.Response, err error, selectedHeaders []string) map[string]string {
	var res = make(map[string]string)

	if err != nil {
		res["Error"] = err.Error()
	} else {
		for key, val := range resp.Header {
			for _, shn := range selectedHeaders {
				if strings.EqualFold(key, shn) {
					res[key] = strings.Join(val, ",")
				}
			}
		}

		res["StatusCode"] = fmt.Sprint(resp.StatusCode)
	}

	return res
}

func printHeadersFor(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	selectedHeaders := []string{"Cache-Control", "Via", "X-Cache"}

	var resp, error = getUrlResponse(url)

	var headers = getInfoAndSelelectedHeaders(resp, error, selectedHeaders)

	for key, value := range headers {
		fmt.Printf("%s\t%s=%s\n", url, key, value)
	}
}

func main() {
	var wg sync.WaitGroup

	urls := getUrlsFromBody(getUrlResponse(os.Args[1]))
	for _, url := range urls {
		wg.Add(1)
		go printHeadersFor(&wg, url)
	}

	wg.Wait()
}
