package app

import (
	"github.com/Grishameister/2020_2_Eternity/server"
	"log"
)

func main() {
	if conn := config.Db; conn == nil {
		log.Fatal("Connection refused")
		return
	}

	srv := server.New(config.Conf)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
