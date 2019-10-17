package main

import (
	"biubiubiu/core"
	"biubiubiu/plugins"
	"github.com/fvbock/endless"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)


func main() {

	proxyMap := core.GetProxyMap()
	configContext := core.GetConfigContext()

	r := chi.NewRouter()
	r.Use(plugins.Cors)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/p/{app}/{uri}", func(w http.ResponseWriter, r *http.Request) {
		app := chi.URLParam(r, "app")
		uri := chi.URLParam(r, "uri")

		if _, ok := proxyMap[app]; !ok {
			log.Println("no proxy server config for " + app)
			http.NotFound(w, r)
			return
		}
		instance := proxyMap[app]
		handler := &core.CommonHandler{Inst: instance, ConfigContext: configContext}
		handler.ServeHTTP(uri, w, r)
	})

	err := endless.ListenAndServe(":3000", r)

	log.Fatal(err)
}