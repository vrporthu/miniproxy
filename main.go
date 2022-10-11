package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strconv"
)

var config Config

func proxyPass(res http.ResponseWriter, req *http.Request) {
	url, _ := url2.Parse(config.ServerUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = rewrite
	proxy.ServeHTTP(res, req)
}

func rewrite(res *http.Response) error {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	for f, t := range config.Replacements {
		b = bytes.Replace(b, []byte(f), []byte(t), -1)
	}
	body := io.NopCloser(bytes.NewReader(b))
	res.Body = body
	res.ContentLength = int64(len(b))
	res.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
		return
	}
	config = c
	http.HandleFunc("/", proxyPass)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
