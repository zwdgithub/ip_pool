package main

import (
	"net/http"
	"time"
)

func main() {

	proxy := NewProxy()
	ticker := time.NewTicker(time.Second * 3)
	kdlGetter := &KDLProxyGetter{}
	go func() {
		for range ticker.C {
			proxy.GetProxyPushTo(kdlGetter)
		}
	}()
	ticker1 := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker1.C {
			proxy.ValidRepeatCheck()
		}
	}()
	go proxy.Run(proxy.redis.GetValidWaitingIp)
	go proxy.Run(proxy.redis.GetWaitingIp)

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		ip, err := proxy.GetProxy()
		if err != nil {
			w.Write([]byte(""))
			return
		}
		w.Write([]byte(ip))
	})
	http.ListenAndServe(":8091", nil)
}
