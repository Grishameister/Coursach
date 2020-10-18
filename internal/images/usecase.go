package images

import "context"

type IUsecase interface {
	PostFrame(c context.Context, bytes []byte) error
}
