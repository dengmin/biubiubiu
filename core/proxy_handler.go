package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const  CacheHeader = "X-Cache"

type CommonHandler struct {
	Inst Instance
	ConfigContext ProxyConfigContext
}


func (handler *CommonHandler)ServeHTTP(uri string, w http.ResponseWriter, r *http.Request){

	instance := handler.Inst

	realIpAddr := realIP(r)

	if len(instance.WhiteIps) > 0{
		ipList := strings.Split(instance.WhiteIps, ",")
		//请求的ip 不在白名单中
		if !constantIp(ipList, realIpAddr) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}
	}

	server, err := instance.GetLoadBalance().Target()
	if err != nil{
		log.Fatal(err)
	}

	remote, err := url.Parse("http://" + server)

	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	r.URL.Scheme = remote.Scheme
	r.URL.Path = uri
	r.Header.Set("Host", instance.Domain)

	rw := newResponseStreamer(w)
	rdr, err := rw.Stream.NextReader()
	if err != nil {
		w.Header().Set(CacheHeader, "SKIP")
		proxy.ServeHTTP(w, r)
		return
	}

	rw.Header().Set(CacheHeader, "MISS")
	go func() {
		proxy.ServeHTTP(rw, r)
		rw.Stream.Close()
	}()
	rw.WaitHeaders()

	b, err := ioutil.ReadAll(rdr)
	if err := rdr.Close(); err != nil{
		log.Print("stream close error" + err.Error())
	}
	fmt.Println(string(b))
	//proxy.ServeHTTP(w, r)
}


func realIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func constantIp(ips []string, realIp string) bool {
	for i := 0; i < len(ips); i++ {
		if ips[i] == realIp {
			return true
		}
	}
	return false
}

