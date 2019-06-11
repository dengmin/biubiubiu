package core

import "biubiubiu/balance"
import "time"

type ProxyConfigContext struct {
	Services []Instance
	LogPretty bool `yaml:"logPretty"`
}

type Instance struct {
	Name string `yaml:"name"`
	Server []string `yaml:"server"`
	Host string `yaml:"host"`
	EnableRateLimit bool `yaml:"enableRateLimit"`
	EnableCache bool `yaml:"enableCache"`
	LoadBalance string `yaml:"loadBalance"`
}

func (inst *Instance) GetLoadBalance() balance.LoadBalance {
	switch inst.LoadBalance {
	case "random":
		return balance.NewRandom(inst.Server, time.Now().UnixNano())
	case "roundRobin":
		return balance.NewRoundRobin(inst.Server)
	default:
		return balance.NewRandom(inst.Server, time.Now().UnixNano())
	}
}