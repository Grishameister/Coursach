package images

import "context"

type IRepository interface {
	PostFrame(c context.Context, bytes []byte) error
}
