package database

import "time"

type DBInterface interface {
	Open() (err error)

	Close() error

	PostFrame(bytes []byte) error
	GetFrameByDate(date time.Time) []byte
	GetLastFrame() []byte
}
