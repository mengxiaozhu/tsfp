package main

import (
	"log"
	"github.com/BurntSushi/toml"
	"github.com/mengxiaozhu/tsfp/proxy"
)

func main() {
	conf := &proxy.Config{}
	if _, err := toml.DecodeFile(".tsfp.toml", conf); err != nil {
		log.Fatal(err.Error())
	}
	// 启动代理
	proxy.NewProxy(conf)
}
