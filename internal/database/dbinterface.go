package database

import "context"

type DBInterface interface {
	Open() (err error)

	Close() error

	PostFrame(c context.Context, bytes []byte) error
}
