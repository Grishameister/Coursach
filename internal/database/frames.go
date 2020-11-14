package database

import (
	"context"
	"log"
	"time"
)

func (db *DB) PostFrame(bytes []byte) error{
	if _, err := db.dbPool.Exec(context.Background(), "insert into frames(bytes, reg_date) values ($1, $2)", bytes, time.Now()); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
