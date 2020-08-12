package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main() {

	proxy := NewProxy()
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for range ticker.C {
			proxy.getProxyPushToRedis()
		}
	}()
	ticker1 := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker1.C {
			proxy.ValidRepeatCheck()
		}
	}()
	go proxy.Run()

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		ip, err := proxy.GetProxy()
		u := "http://myip.ipip.net/"
		if err != nil {
			w.Write([]byte(""))
			return
		}
		p := func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://" + ip)
		}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           p,
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Second * 5,
		}
		resp, err := client.Get(u)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(ip + "\n" + string(bytes)))
	})
	http.ListenAndServe(":80", nil)
}
