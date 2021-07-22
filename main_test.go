package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
	"testing"
)

var selectedHeaders = []string{"Cache-Control", "Via", "X-Cache"}

func TestGetInfoAndSelelectedHeadersWithError(t *testing.T) {

	info := getInfoAndSelelectedHeaders(nil, errors.New("error"), selectedHeaders)

	assertNotNil(info)
	assertEqualFold(info["Error"], "error")
}

func TestGetInfoAndSelelectedHeaders200(t *testing.T) {

	var httpHeaders = map[string][]string{}
	httpHeaders["Cache-Control"] = []string{"value1", "value2"}
	httpHeaders["Via"] = []string{"value1"}
	httpHeaders["X-Cache"] = []string{"value1"}

	var resp = createHttpResponse(200, httpHeaders)

	info := getInfoAndSelelectedHeaders(&resp, nil, selectedHeaders)

	assertNotNil(info)

	assertEqualFold(info["StatusCode"], "200")
	assertEqualFold(info["Cache-Control"], "value1,value2")
	assertEqualFold(info["Via"], "value1")
	assertEqualFold(info["X-Cache"], "value1")

}

func TestGetInfoAndSelelectedHeaders404(t *testing.T) {

	var resp = createHttpResponse(404, map[string][]string{})

	info := getInfoAndSelelectedHeaders(&resp, nil, selectedHeaders)

	assertNotNil(info)

	assertEqualFold(info["StatusCode"], "404")
}

func TestHasHttpPrefix(t *testing.T) {

	assertTrue(hasHttpPrefix("https://plweb.se"))
	assertTrue(!hasHttpPrefix("#menu"))
}

func createHttpResponse(statusCode int, headers map[string][]string) http.Response {

	resp := http.Response{
		Status:           "",
		StatusCode:       statusCode,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           headers,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          map[string][]string{},
		Request:          &http.Request{},
		TLS:              &tls.ConnectionState{},
	}

	return resp
}

func assertTrue(b1 bool) {
	if b1 != true {
		log.Fatalf("%v != true", b1)
	}
}

func assertEqualFold(s1 string, s2 string) {
	if !strings.EqualFold(s1, s2) {
		log.Fatalf("%s != %s", s1, s2)
	}
}

func assertNotNil(o interface{}) {
	if o == nil {
		log.Fatalf("%s == nil", o)
	}
}
