package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	checkUrl = "https://www.alipay.com/"
)

type ProxyProcess struct {
	redis *RedisUtil
}

func NewProxy() *ProxyProcess {
	return &ProxyProcess{redis: NewRedis()}
}

func (proxy *ProxyProcess) Valid(ip string, checkChan chan bool) bool {
	log.Printf("start check %s", ip)
	result := true
	defer func() {
		<-checkChan
		if err := recover(); err != nil {
			log.Printf("recover err: %v", err)
			result = false
		}
	}()
	// TODO check
	p := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + ip)
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           p,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 3,
	}
	start := time.Now().Unix()
	_, err := client.Get(checkUrl)
	if err != nil {
		proxy.redis.DeleteProxy(ip)
		log.Printf("use proxy get checkUrl error is: %v", err)
		return false
	}
	proxy.redis.AddProxy(ip)
	log.Printf("get checkUrl %s success, cost %d s", checkUrl, time.Now().Unix()-start)
	return result
}

func (proxy *ProxyProcess) GetProxyPushTo(getter ProxyGetter) {
	list, err := getter.GetProxy()
	if err != nil {
		log.Printf("getProxyPushToRedis get checkUrl err: %v", err)
	}
	for _, ip := range list {
		log.Printf("push %s to waiting list", ip)
		proxy.redis.AddIpToWaitingList(ip)
	}
}

func (proxy *ProxyProcess) ValidRepeatCheck() {
	list, err := proxy.redis.AllValidIp()
	if err != nil {
		log.Printf("get all valid ip error: %v", err)
		return
	}
	for _, ip := range list {
		proxy.redis.AddValidIpToWaitingList(ip)
	}
}

func (proxy *ProxyProcess) Run(f func() (string, error)) {
	var count int64 = 0
	checkChan := make(chan bool, 100)
	for {
		atomic.AddInt64(&count, 1)
		log.Printf("the %d times", count)
		ip, err := f()
		if err != nil {
			log.Printf("get waiting ip error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		checkChan <- true
		go proxy.Valid(ip, checkChan)
	}
}

func (proxy *ProxyProcess) GetProxy() (string, error) {
	ip, err := proxy.redis.GetProxy()
	if err != nil {
		return ip, err
	}
	proxy.redis.AddValidIpToWaitingList(ip)
	return ip, err
}
