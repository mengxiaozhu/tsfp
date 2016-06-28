package main

import (
	"github.com/mengxiaozhu/tsfp"
	"github.com/BurntSushi/toml"
	"log"
)

func main() {
	conf := &tsfp.Config{}
	if _, err := toml.DecodeFile(".tsfp.toml", conf); err != nil {
		log.Fatal(err.Error())
	}
	// 启动代理
	tsfp.NewProxy(conf)
}
