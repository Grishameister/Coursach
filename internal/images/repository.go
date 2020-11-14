package images

type IRepository interface {
	PostFrame(bytes []byte) error
}
