package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//读取配置文件
func NewProxyConfigMap() (map[string]Instance, ProxyConfigContext, error){
	proxyMap := make(map[string]Instance)
	configContext := new(ProxyConfigContext)
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil{
		log.Panic("read config file error", err)
	}
	err = yaml.Unmarshal(yamlFile, configContext)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	for i := 0; i < len(configContext.Services); i++ {
		instance := configContext.Services[i]
		proxyMap[instance.Name] = instance
	}
	return proxyMap, *configContext, nil
}


