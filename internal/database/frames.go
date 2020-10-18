package database

import (
	"context"
	"log"
	"time"
)

func (db *DB) PostFrame(c context.Context, bytes []byte) error{
	if _, err := db.dbPool.Exec(c, "insert into frames(bytes, reg_date) values ($1, $2)", bytes, time.Now()); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
