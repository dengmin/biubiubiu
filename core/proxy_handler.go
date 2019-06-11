package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	Inst Instance
	ConfigContext ProxyConfigContext
}


func (handler *Handler)ServeHTTP(uri string, w http.ResponseWriter, r *http.Request){

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

	proxy.ModifyResponse = rewriteBody
	proxy.ServeHTTP(w, r)
}

func rewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  err
	}
	//log.Println(string(b))
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
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
