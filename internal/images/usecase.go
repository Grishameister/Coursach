package images

type IUsecase interface {
	PostFrame(bytes []byte) error
}
