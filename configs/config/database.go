package config

import (
	"github.com/jackc/pgx/pgxpool"
)

type DB struct {
	dbPool *pgxpool.Conn
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Open() error {

}
