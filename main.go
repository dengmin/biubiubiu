package main

import (
	"biubiubiu/core"
	"biubiubiu/plugins"
	"bytes"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func rewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  err
	}
	log.Println(string(b))
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	//////test := string(b)
	//////log.Println(test)
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}


func main() {

	configMap, configContext, _ := core.NewProxyConfigMap()

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		PrettyPrint: configContext.LogPretty,
	}

	r := chi.NewRouter()
	r.Use(plugins.NewStructuredLogger(logger))
	r.Use(plugins.Cors)



	r.HandleFunc("/p/{app}/{uri}", func(w http.ResponseWriter, r *http.Request) {
		app := chi.URLParam(r, "app")
		uri := chi.URLParam(r, "uri")

		if _, ok := configMap[app]; !ok {
			log.Println("no proxy server config for " + app)
			http.NotFound(w, r)
			return
		}
		instance := configMap[app]

		server, err := instance.GetLoadBalance().Target()
		log.Println("get Server addr:" + server)
		if err != nil{
			log.Fatal(err)
		}

		u, err := url.Parse("http://" + server)

		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		r.URL.Scheme = u.Scheme
		r.URL.Path = uri
		r.Header.Set("Host", instance.Host)
		proxy.ModifyResponse = rewriteBody
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":3000", r))
}