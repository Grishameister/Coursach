package database

type DBInterface interface {
	Open() (err error)

	Close() error

	PostFrame(bytes []byte) error
}
