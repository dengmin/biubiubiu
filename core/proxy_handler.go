package core

import (
	"biubiubiu/cache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const  CacheHeader = "X-Cache"

type CommonHandler struct {
	Inst Instance
	ConfigContext ProxyConfigContext
}

var cacheMap = make(map[string]*cache.Cache)

func (handler *CommonHandler)ServeHTTP(uri string, w http.ResponseWriter, r *http.Request){

	instance := handler.Inst

	realIpAddr := getRealIP(r)

	if len(instance.WhiteIps) > 0{
		ipList := strings.Split(instance.WhiteIps, ",")
		//请求的ip 不在白名单中
		if !constantIp(ipList, realIpAddr) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}
	}

	var cacheKey = ""
	var c *cache.Cache

	if instance.EnableCache {

		if _, ok := cacheMap[instance.Name]; !ok {
			c = cache.New(5*time.Minute, 10*time.Minute)
			cacheMap[instance.Name] = c
		}else{
			c = cacheMap[instance.Name]
		}

		//清空缓存
		if uri == "_clean_cache" {
			c.Flush()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("cache cleaned!"))
			return
		}

		cacheKey := handler.buildKeys(r)

		//从缓存中读取数据写入到response
		v, found := c.Get(cacheKey)
		if found {
			log.Println("get entry from cache")
			w.Header().Set(CacheHeader, "HIT")
			w.Write([]byte(v.(string)))
			return
		}
		w.Header().Set(CacheHeader, "MISS")
	}

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

	rw := newResponseStreamer(w)
	rdr, err := rw.Stream.NextReader()
	if err != nil {
		proxy.ServeHTTP(w, r)
		return
	}

	go func() {
		proxy.ServeHTTP(rw, r)
		_ = rw.Stream.Close()
	}()
	rw.WaitHeaders()

	b, err := ioutil.ReadAll(rdr)
	if err := rdr.Close(); err != nil{
		log.Print("stream close error" + err.Error())
	}

	//只对GET请求缓存
	if instance.EnableCache && r.Method == "GET" {
		go func() {
			c.Set(cacheKey, string(b), cache.DefaultExpiration)
		}()
	}

	//状态码500报警
	if rw.StatusCode > 500{
		//TODO
	}

}

//通过参数生成缓存的key,并md5
func (handler *CommonHandler) buildKeys(r *http.Request) string{
	//cacheKeyFormat := handler.Inst.CacheKey
	//log.Println(cacheKeyFormat)
	//获取参数
	t := r.URL.Query().Encode()
	realIp := getRealIP(r)
	target := realIp+":"+r.URL.Path+":"+t
	return md5Encrypt(target)
}
