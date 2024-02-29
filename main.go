package main

import (
	"server/conf"
	"server/server"
)

func main() {
	conf.InitConfig("dev")
	config := conf.GetConfig()
	server, err := server.NewServer(config)
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
