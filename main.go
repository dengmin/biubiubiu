package main

import (
	"biubiubiu/core"
	"biubiubiu/plugins"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)


func main() {

	proxyMap, configContext, _ := core.NewProxyConfigMap()

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

		if _, ok := proxyMap[app]; !ok {
			log.Println("no proxy server config for " + app)
			http.NotFound(w, r)
			return
		}
		instance := proxyMap[app]

		handler := &core.Handler{Inst: instance, ConfigContext:configContext}
		handler.ServeHTTP(uri, w, r)
	})

	log.Fatal(http.ListenAndServe(":3000", r))
}