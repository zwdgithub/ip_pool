package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

var (
	r        = rand.New(rand.NewSource(time.Now().Unix()))
	checkUrl = "http://myip.ipip.net/"
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

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func (proxy *ProxyProcess) getProxyPushToRedis() {
	client := &http.Client{}
	resp, err := client.Get("https://dev.kdlapi.com/api/getproxy/?orderid=949722172204228&num=100&protocol=2&method=1&an_ha=1&sp1=1&sp2=1&quality=1&sort=2&sep=2")
	if err != nil {
		log.Printf("getProxyPushToRedis get checkUrl err: %v", err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	for _, ip := range strings.Split(string(bytes), "\n") {
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
		proxy.redis.AddIpToWaitingList(ip)
	}
}

func (proxy *ProxyProcess) Run() {
	var count int64 = 0
	checkChan := make(chan bool, 100)
	for {
		atomic.AddInt64(&count, 1)
		log.Printf("the %d times", count)
		ip, err := proxy.redis.GetWaitingIp()
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
	return proxy.redis.GetProxy()
}
