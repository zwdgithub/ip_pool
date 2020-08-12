package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	r := NewRedis()
	v, err := r.GetWaitingIp()
	t.Log(v)
	t.Log(err)
}

func TestHttp(t *testing.T) {
	ip := "171.35.142.60:9999"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + ip)
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           proxy,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 3,
	}
	_, err := client.Get(checkUrl)
	fmt.Println(err)

}
