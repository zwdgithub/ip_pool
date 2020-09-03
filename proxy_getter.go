package main

import (
	"fmt"
	xhttp "github.com/zwdgithub/simple_http"
	"log"
	"math/rand"
	"strings"
)

type ProxyGetter interface {
	// 获取代理
	GetProxy() ([]string, error)
}

// 快代理
type KDLProxyGetter struct {
}

func (getter *KDLProxyGetter) GetProxy() ([]string, error) {
	result := make([]string, 0)
	n := rand.Intn(99999)
	var carrier int
	switch n % 3 {
	case 0:
		carrier = 1
	case 1:
		carrier = 2
	case 3:
		carrier = 6
	}
	url := fmt.Sprintf("http://dev.kdlapi.com/api/getproxy/?orderid=949722172204228&num=100&carrier=%d&protocol=2&method=1&an_ha=1&sep=2", carrier)
	content, err := xhttp.NewHttpUtil().Get(url).RContent()
	if err != nil {
		log.Printf("KDLProxyGetter get proxy error, err: %v", err)
		return result, err
	}
	for _, ip := range strings.Split(content, "\n") {
		log.Printf("push %s to waiting list", ip)
		result = append(result, ip)
	}
	return result, nil
}
