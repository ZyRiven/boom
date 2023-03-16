package main

import (
	"gohello/config"
	"gohello/duang"
	"gohello/route"
)

func main() {
	r := duang.New()
	route.Init(r)
	duang.WebSocketMain() // webSocket,不用可以删除
	r.Run(config.RetServer().Address)
}
