package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"testing"
)

func TestGetInfoAndSelelectedHeaders(t *testing.T) {
	var headers = []string{"test", "test2"}

	info := getInfoAndSelelectedHeaders(nil, errors.New("error"), headers)

	if info == nil {
		log.Fatal()
	}

	if info["error"] != "error" {
		log.Fatal()
	}

}

func TestGetInfoAndSelelectedHeaders2(t *testing.T) {
	var headers = []string{"test", "test2"}

	resp := http.Response{
		Status:           "",
		StatusCode:       0,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           map[string][]string{},
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          map[string][]string{},
		Request:          &http.Request{},
		TLS:              &tls.ConnectionState{},
	}

	info := getInfoAndSelelectedHeaders(&resp, nil, headers)

	if info == nil {
		log.Fatal()
	}

	if info["statusCode"] != "0" {
		log.Fatal()
	}

}
