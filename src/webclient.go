package main

import (
	"net/http"
	"strings"
	"time"
)

var wCli *http.Client

func init() {
	wCli = &http.Client{}
	wCli.Timeout = time.Second
}

func Do(method string, url string, body string) bool {
	req, er := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(body))
	if er != nil {
		return false
	}
	req.Header.Set("Content-Type", "applciation/json")
	r, er := wCli.Do(req)
	if er != nil || r.StatusCode != 200 {
		return false
	}
	return true
}
