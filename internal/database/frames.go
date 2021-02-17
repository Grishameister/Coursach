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

func (db *DB) GetFrameByDate(date time.Time) []byte {
	bytes := make([]byte, 0, 0)
	if err := db.dbPool.QueryRow(context.Background(), "select bytes from frames order by $1 <-> reg_date limit 1", date).Scan(&bytes); err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}

func (db *DB) GetLastFrame() []byte {
	bytes := make([]byte, 0, 0)
	if err := db.dbPool.QueryRow(context.Background(), "select bytes from frames order by reg_date desc limit 1").Scan(&bytes); err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}
