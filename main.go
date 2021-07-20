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

func getUrlsFromBody(resp *http.Response, err error) []string {
	defer resp.Body.Close()
	urls := []string{}

	if err != nil {
		log.Print(err)
		return urls
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Print(err)
		return urls
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {

		url, _ := s.Attr("href")
		if strings.HasPrefix(url, "http") {
			urls = append(urls, url)
		}
	})

	return urls
}

func getInfoAndSelelectedHeaders(resp *http.Response, err error, selectedHeaders []string) map[string]string {
	var res = make(map[string]string)

	if err != nil {
		res["error"] = err.Error()
	} else {
		for key, val := range resp.Header {
			for _, shn := range selectedHeaders {
				if strings.EqualFold(key, shn) {
					res[key] = strings.Join(val, ",")
				}
			}
		}

		res["statusCode"] = fmt.Sprint(resp.StatusCode)
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
