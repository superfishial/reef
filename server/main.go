package main

import (
	"github.com/superfishial/reef/server/api"
	"github.com/superfishial/reef/server/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	api.StartServer(config)
}
