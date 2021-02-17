package images

import "time"

type IRepository interface {
	PostFrame(bytes []byte) error

	GetFrameByDate(date time.Time) []byte
	GetLastFrame() []byte
}
