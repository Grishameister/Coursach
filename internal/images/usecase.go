package images

import "time"

type IUsecase interface {
	PostFrame(bytes []byte) error

	GetFrameByDate(date time.Time) []byte
	GetLastFrame() []byte
}
