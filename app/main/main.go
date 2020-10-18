package main

import (
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/database"
	"github.com/Grishameister/Coursach/internal/queue"
	"github.com/Grishameister/Coursach/server"
	"log"
)

func main() {

	q := queue.NewQueue()
	db := database.NewDB(&config.Conf.Db)
	if err := db.Open(); err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	srv := server.New(config.Conf, db, q)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
