package core

import (
	"log"
	"net/http"
	"net/url"
)


type CommonHandler struct {
	Inst Instance
	ConfigContext ProxyConfigContext
}


func (handler *CommonHandler)ServeHTTP(uri string, w http.ResponseWriter, r *http.Request){

	instance := handler.Inst

	server, err := instance.GetLoadBalance().Target()
	if err != nil{
		log.Fatal(err)
	}

	remote, err := url.Parse("http://" + server)

	if err != nil {
		log.Fatal(err)
	}

	proxy := NewSingleHostReverseProxy(remote)
	r.URL.Scheme = remote.Scheme
	r.URL.Path = uri
	r.Header.Set("Host", instance.Domain)


	go func() {
		proxy.ServeHTTP(w, r)
	}()


}

