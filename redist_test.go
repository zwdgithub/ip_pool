package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
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

func Get(ctx context.Context, c chan string) {
	n := rand.Intn(10)
	time.Sleep(n)
	for {
		select {
		case <-ctx.Done():

		}
	}
}

func TestGet(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	proxy := NewProxy()
	ip, err := proxy.GetProxy()
	go Get(ctx, 1)
}
