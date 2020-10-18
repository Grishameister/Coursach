package app

import (
	"github.com/Grishameister/2020_2_Eternity/configs/config"
	"github.com/Grishameister/2020_2_Eternity/internal/app/server"
	"log"
)

func main() {
	if conn := config.Db; conn == nil {
		log.Fatal("Connection refused")
		return
	}

	defer Close()
	srv := server.New(config.Conf)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
