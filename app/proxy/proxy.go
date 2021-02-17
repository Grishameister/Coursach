package main

import (
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/serverProxy"
	"log"
)

func main() {
	srv := serverProxy.New(config.Conf)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
